# Plan de pruebas y cobertura de casos de uso

Este documento planifica el trabajo de testing de la librería: qué falta, qué
bugs ya salieron a la luz al analizarla, y en qué orden atacarlos.

El objetivo no es subir el porcentaje de cobertura: es que ninguna combinación
razonable de entradas produzca un resultado silenciosamente incorrecto, un
pánico, o un error donde debería haber un número.

---

## 1. Estado actual (medido, no estimado)

Cobertura por paquete, con `go test -cover ./...`:

| Paquete | Cobertura | Tests | Funciones exportadas | Tests / API |
| --- | ---: | ---: | ---: | ---: |
| `finance/gradients` | **61.3 %** | 11 | 47 | 0.23 |
| `finance/investment` | **63.6 %** | 18 | 30 | 0.60 |
| `finance/tvm` | **73.3 %** | 11 | 29 | 0.38 |
| `finance/bonds` | **74.3 %** | 9 | 28 | 0.32 |
| `finance/depreciation` | **80.2 %** | 8 | 18 | 0.44 |
| `finance/daycount` | 84.4 % | 14 | 17 | 0.82 |
| `finance/annuities` | 84.7 % | 90 | 149 | 0.60 |
| `finance/returns` | 85.4 % | 42 | 78 | 0.54 |
| `finance/compoundinterest` | 85.9 % | 152 | 213 | 0.71 |
| `charts` (módulo aparte) | 86.3 % | 6 | — | — |
| `finance/loans` | 90.1 % | 26 | 60 | 0.43 |
| `money` | 90.9 % | 82 | 139 | 0.59 |
| `finance/term` | 95.7 % | 4 | 8 | 0.50 |
| `decimal` | 96.6 % | 156 | 230 | 0.68 |
| `finance/simpleinterest` | 97.0 % | 29 | 59 | 0.49 |
| **Total** | **85.2 %** | 632 | ~1 105 | |

Otros datos del inventario:

- **39 funciones con 0 % de cobertura.** Entre ellas la ruta completa de
  respaldo de los solucionadores: `irrBisection`, `bisect`, `irrCandidates`,
  `xirrBisection`, `xbisect`. Es decir: el camino que se ejecuta *precisamente
  cuando Newton-Raphson falla* nunca se ha ejecutado en una prueba.
- **0 tests de tipo `Example`.** Para una librería pública esto es una pérdida
  doble: no hay ejemplos ejecutables en godoc y los fragmentos de código de los
  comentarios no se compilan, así que pueden quedar desactualizados sin que
  nadie se entere.
- **0 objetivos de fuzzing**, pese a que hay un motor decimal de 128 bits y
  varios parsers (`FromString`, JSON, SQL) que son candidatos naturales.
- **Benchmarks sólo en 3 paquetes de dominio** (`annuities`,
  `compoundinterest`, `simpleinterest`) más `decimal` y `money`.
- El grueso de las pruebas verifica el **camino feliz con un solo juego de
  números**. Casi no hay pruebas de invariantes, de reciprocidad entre
  funciones, ni de coherencia entre paquetes.

---

## 2. Hallazgos: defectos reales detectados durante el análisis

Estos no son riesgos hipotéticos. Cada uno se reprodujo ejecutando la librería.
Son la evidencia de por qué el plan está ordenado como está.

### 2.1 Pánico en funciones que devuelven `error` (tasa 0 %) — crítico

Las cuatro funciones de pago de `finance/annuities` entran en pánico
(`division by zero`) con una tasa de interés del 0 %, en lugar de devolver el
error que su firma promete:

```
PaymentFromPresentValue             -> PANIC: division by zero
PaymentFromFutureValue              -> PANIC: division by zero
AnticipatePaymentFromPresentValue   -> PANIC: division by zero
AnticipatePaymentFromFutureValue    -> PANIC: division by zero
```

Causa: `finance/annuities/root.go` y `deferred.go` usan los ayudantes que
entran en pánico (`MustPow`, `MustDiv`) dentro de funciones que sí devuelven
`error`. Con `i = 0`, el denominador `(1+i)^n − 1` vale cero y `MustDiv`
explota.

Un préstamo o un plan de ahorro al 0 % (financiación promocional, préstamo
entre familiares, aporte sin rendimiento) es un caso de uso corriente, y el
resultado analítico existe y es trivial: `pago = valor / n`.

Hay 39 llamadas a ayudantes `Must*` repartidas por `finance/` y `money/`. Cada
una es un pánico potencial dentro de una API que devuelve errores.

### 2.2 Ceros silenciosos en la conversión de tasas — crítico

`finance/compoundinterest` convierte entre cinco tipos de tasa. La matriz de
conversión (valor 0.12, capitalización mensual) da esto:

| tipo origen | `RatePeriodic` | `RateNominal` | `RateEffectyAnnually` | `RateAnticipateNominal` | `RateAnticipatePeriodic` |
| --- | --- | --- | --- | --- | --- |
| `periodic` | 0.12 | 1.44 | 2.8959… | **0** | **0** |
| `annual` | 0.00948… | 0.11386… | 0.12 | 0.11279… | 0.00939… |
| `nominal` | 0.01 | 0.12 | 0.12682… | **0** | **0** |
| `anticipatePeriodic` | **0** | **0** | **0** | 1.44 | 0.12 |
| `anticipateNominal` | **0** | **0** | **0** | 0.12 | 0.01 |

**12 de 25 combinaciones devuelven `0` con `err == nil`.** Las funciones están
escritas como una cadena de `if` sin `else` ni caso por defecto: cuando el tipo
de tasa no coincide con ninguna rama, la variable de resultado se queda en su
valor cero y se devuelve como si fuera un cálculo válido.

Pasar de una tasa vencida a una anticipada está perfectamente definido
(`d = i/(1+i)`), así que aquí no falta una validación: falta la conversión, y
en su lugar se propaga un 0 hacia todo lo que dependa de esa tasa.

Es el peor tipo de fallo posible en una librería financiera: no rompe nada,
sólo devuelve una cifra equivocada.

### 2.3 Series de gradientes al 0 % — importante

`gradients.Arithmetic.Present()` y `.Future()` devuelven `division by zero` con
una tasa del 0 %. Aquí no hay pánico, pero el límite analítico existe
(`PV = Σ pagos`) y devolver un error obliga al usuario a resolver a mano un
caso legítimo. `Geometric` sí trata su singularidad (`g == i`); `Arithmetic` no
trata la suya.

### 2.4 `term.Daily.MonthsPerPeriod()` es inconsistente consigo mismo — importante

`term.Daily.PeriodsPerYear()` devuelve 365, pero `MonthsPerPeriod()` devuelve
la constante escrita a mano `0.03333333` (es decir, 1/30). El valor coherente
con 365 períodos al año es `12/365 = 0.032876712…`. Un error relativo del 1.4 %
que se cuela en cualquier cálculo diario que mezcle ambas funciones.

### 2.5 Regla de fin de mes en 30/360 — a decidir

`daycount.thirty360Days` documenta que aplica «los ajustes estándar de fin de
mes», pero sólo implementa la regla del día 31. La convención US (NASD)
30/360 también trata febrero: del 29-feb-2024 al 31-ago-2024 la librería
devuelve 182 días, mientras que la convención con regla de febrero da 180.

No está claro que sea un bug — depende de qué variante se quiera ofrecer —
pero el código y la documentación no dicen lo mismo, y ninguna prueba fija el
comportamiento.

### 2.6 `IRR` con varias raíces devuelve una sola, sin avisar — a decidir

Con flujos `[-1000, 6000, -11000, 6000]` (dos cambios de signo, varias TIR
matemáticamente válidas), `IRR` devuelve ~0 sin ninguna señal de que la
solución no es única. Es la limitación clásica de la TIR y es aceptable, pero
debe quedar documentada y fijada por una prueba, no descubierta por un usuario
en producción.

### 2.7 `BuildSchedule` trunca los períodos fraccionarios en silencio

`annuities.BuildSchedule(..., nper = 5.7)` genera 5 períodos sin error. O se
rechaza (`ErrInvalidPeriods`) o se documenta el truncamiento; hoy no se hace
ninguna de las dos cosas.

---

## 3. Estrategia: seis tipos de prueba, no uno

Las pruebas actuales son casi todas del tipo 1. Los defectos de la sección 2 se
encuentran con los tipos 2, 3 y 4.

**1. Tabla de casos (lo que ya hay).** Entrada conocida → salida esperada.
Sigue siendo la base; hay que ampliarla a los valores frontera de cada función:
cero, negativo, un solo período, tasas extremas, monedas sin decimales.

**2. Invariantes.** Propiedades que deben cumplirse para *cualquier* entrada
válida, no para un número concreto:

- Un cuadro de amortización cierra en saldo 0 y la suma de capital equivale al
  principal.
- La suma de una depreciación es exactamente `costo − residual`.
- `Allocate` siempre suma el monto original, sin importar los ratios.
- `NPV(IRR(flujos), flujos) ≈ 0`.
- `Price(YTM(precio)) ≈ precio`.
- La duración de un bono cupón cero es igual a su plazo.

**3. Reciprocidad (ida y vuelta).** Cada par de funciones inversas debe cerrar
el círculo: `PV → PMT → PV`, `nominal → efectiva → nominal`, `vencida →
anticipada → vencida`, `Money → JSON → Money`, `decimal → string → decimal`.
La sección 2.2 aparece sola en cuanto se escribe la prueba de ida y vuelta de
tasas para las 25 combinaciones.

**4. Barrido de la matriz de la API.** Cuando una función acepta un tipo
enumerado (tipo de tasa, frecuencia, convención, moneda), la prueba recorre
*todos* los valores del enumerado, no uno. Los 12 ceros silenciosos estaban en
las celdas que ninguna prueba visitaba.

**5. Ejemplos ejecutables (`Example…`).** Documentación que el compilador
verifica. Un `Example` por concepto principal de cada paquete; los fragmentos
que hoy viven en comentarios se convierten en ejemplos reales.

**6. Fuzzing.** Para las fronteras que reciben texto o bytes arbitrarios:
`decimal.FromString`, `UnmarshalJSON`, `Scan`/`Value` de SQL, y las rutinas
aritméticas de 128/256 bits contra `math/big` como oráculo.

---

## 4. Plan por paquete

Prioridad: **P0** desbloquea correcciones de bugs conocidos, **P1** es riesgo
alto sin cobertura, **P2** es consolidación.

### `finance/compoundinterest` — P0

El paquete del que dependen `annuities`, `gradients` y buena parte del resto.
Un error aquí contamina todo lo demás.

- Matriz de conversión completa: 5 tipos de tasa × 5 conversiones × 7
  frecuencias. Cada celda debe dar un valor correcto o un error explícito —
  nunca un cero silencioso (§2.2).
- Ida y vuelta de conversiones: `X → Y → X` recupera el valor inicial dentro de
  la tolerancia.
- Equivalencia financiera: una tasa nominal del 12 % capitalizable mensualmente
  y su efectiva anual del 12.6825 % deben producir el mismo valor futuro.
- Coherencia anticipada/vencida: `d = i/(1+i)` verificada en ambos sentidos.
- Frontera: tasa 0 %, tasa negativa (tipos de interés negativos existen), tasa
  anticipada ≥ 100 % (degenerada: debe dar error), 1 período, períodos
  fraccionarios.
- `GetEqualsRateInterestPeriods` con cada combinación de frecuencia de tasa y
  de período (mensual con período anual, etc.).

### `finance/annuities` — P0

- Tasa 0 % en las cuatro funciones de pago, las diferidas, las de valor
  presente/futuro y las de número de períodos (§2.1). Resultado esperado
  analítico, no pánico.
- Auditar las 39 llamadas a `Must*` en código de librería: cada una dentro de
  una función que devuelve `error` debe convertirse en propagación de error.
  La prueba que lo fija: recorrer la API pública con entradas degeneradas y
  afirmar que **ninguna** llamada entra en pánico.
- Invariantes del cuadro de amortización: saldo final 0, Σ capital = principal,
  Σ interés = Σ pago − principal, saldo monótono decreciente.
- Reciprocidad ordinaria/anticipada: `pago_anticipado = pago_vencido / (1+i)`.
- `BuildSchedule` con `nper` fraccionario, negativo, cero y muy grande (§2.7).
- `WriteCSVTo` con un `io.Writer` que falla, para cubrir la ruta de error.

### `finance/gradients` — P1 (cobertura más baja: 61.3 %)

- Tasa 0 % en `Arithmetic.Present`/`Future` (§2.3).
- `Geometric` con `g == i` (ya tratado en el código, sin prueba), `g > i`,
  `g < 0`, `i < 0`.
- Constructores de `builder.go`: `AnnualRate`, `Monthly`, `Quarterly`,
  `MustBuild`, `MustPresent` están a 0 % en ambos constructores.
- Descuadre de monedas entre `firstPayment` y `gradient`.
- Invariante contra `annuities`: un gradiente con `gradient = 0` debe dar
  exactamente lo mismo que una anualidad ordinaria; un geométrico con
  `g = 0`, igual.
- Invariante contra la suma directa: para n pequeño, el VP calculado debe
  coincidir con descontar los pagos uno a uno.

### `finance/investment` — P1 (rutas de respaldo a 0 %)

- **Forzar la bisección.** Construir flujos donde Newton-Raphson diverja o se
  salga del dominio (−1, ∞) para ejecutar `irrBisection`, `bisect`,
  `irrCandidates`, `xirrBisection`, `xbisect`.
- TIR múltiple: fijar el comportamiento documentado (§2.6).
- TIR sin raíz real → `ErrNoConvergence`; sin cambio de signo →
  `ErrNoSignChange`; flujos vacíos, monedas mezcladas.
- TIR muy negativa (cerca de −99 %) y muy alta (> 1000 %), en los extremos de
  la lista de candidatos.
- `XIRR`/`XNPV`: fechas desordenadas, fechas repetidas, un solo flujo, años
  bisiestos, huecos de más de un año.
- Invariante: `NPV(IRR(f), f) ≈ 0` y `XNPV(XIRR(f), f) ≈ 0` sobre un conjunto
  variado de flujos.
- Coherencia: para flujos con fechas separadas exactamente un año, `XIRR` debe
  coincidir con `IRR`.
- Perpetuidades: `g ≥ r` (divergente, debe dar error), `r = 0`, `g` negativa.

### `finance/bonds` — P1 (9 tests para 28 funciones exportadas)

- Precio: a la par (cupón = rendimiento → precio = valor nominal), con prima,
  con descuento, cupón cero, rendimiento negativo, un solo período restante.
- Ida y vuelta `Price` ↔ `YTM` en una rejilla de cupones, plazos y frecuencias.
- Duración: la Macaulay de un cupón cero es su plazo (verificado: correcto —
  conviene fijarlo con una prueba); duración modificada = Macaulay/(1+y/f);
  duración menor que el plazo para bonos con cupón.
- Convexidad positiva; la aproximación de segundo orden con duración y
  convexidad debe acercarse al cambio real de precio ante un desplazamiento
  pequeño del rendimiento.
- `MustMacaulayDuration`, `MustModifiedDuration`, `MustConvexity` están a 0 %:
  probar que devuelven el valor y que entran en pánico ante términos inválidos.
- `AccruedInterest`: liquidación en la fecha del cupón (0), en la fecha del
  siguiente cupón (cupón completo), fuera del período, con las cuatro
  convenciones de conteo de días.
- Errores: frecuencia 0, períodos 0, rendimiento ≤ −frecuencia, precio ≤ 0.

### `finance/tvm` — P1

- Las cinco incógnitas (`FV`, `PV`, `PMT`, `N`, `Rate`) resueltas en círculo:
  fijar cuatro variables, resolver la quinta, volver a resolver la primera.
- `Ordinary` (0 % de cobertura) y anualidades anticipadas: `debido` frente a
  `vencido` en las cinco funciones.
- Tasa 0 % en las cinco (las fórmulas se degeneran).
- `SolveN` cuando no hay solución (el pago no cubre ni el interés) →
  comportamiento definido y probado.
- `SolveRate`: forzar `bisectRate`, señales de no convergencia, tasas
  negativas.
- `MustSolveRate`, `MustSolvePV`, `MustSolveN` (0 %): valor y pánico.
- Coherencia con `annuities` y `loans`: el mismo préstamo resuelto con
  `tvm.SolvePMT` y con `annuities.PaymentFromPresentValue` debe dar la misma
  cuota.

### `finance/depreciation` — P1

- Invariante en los cuatro métodos: `Σ depreciación == costo − residual` y
  `valor en libros final == residual`. (Verificado en línea recta, DDB y SYD;
  falta fijarlo con pruebas y comprobar MACRS.)
- El valor en libros nunca cae por debajo del residual en ningún año
  intermedio.
- Saldo decreciente: el año exacto del cambio a línea recta en DDB.
- MACRS: todos los períodos de recuperación soportados, la convención de medio
  año, y `ErrUnsupportedRecovery` para los no soportados.
- `MustDecliningBalance`, `MustDoubleDecliningBalance`, `MustSumOfYearsDigits`
  están a 0 %.
- Frontera: vida útil de 1 año, residual = 0, residual = costo, residual >
  costo (error, ya verificado), monedas distintas entre costo y residual.

### `finance/daycount` y `finance/term` — P1

- Decidir y fijar la regla de febrero en 30/360 (§2.5): o se implementa la
  variante US completa, o se cambia la documentación y se prueba el
  comportamiento actual.
- Corregir `Daily.MonthsPerPeriod` (§2.4) y añadir la prueba de coherencia
  interna: `MonthsPerPeriod × PeriodsPerYear == 12` para toda frecuencia.
- Las cuatro convenciones × períodos que cruzan años bisiestos, años completos,
  el mismo día, fin de mes, 29 de febrero.
- Actual/Actual ISDA sobre períodos de varios años que empiezan y terminan en
  años bisiestos.
- Entradas con zona horaria y hora del día distintas de medianoche, para
  confirmar que `normalize` las neutraliza.
- `Convention.String()` (50 % de cobertura) para todos los valores, incluido
  uno inválido.

### `money` — P2

- `Allocate`: ya se comprobó que suma correctamente con montos negativos y con
  residuos; falta fijarlo. Añadir: residuo mayor que el número de partes,
  ratios con ceros intercalados, una sola parte, monedas sin decimales (JPY,
  verificado: 1000/3 → 334+333+333) y de tres decimales (BHD, KWD).
- `MulDecimal`, `DivDecimal`, `MustDivDecimal`, `RoundBank` están a 0 % pese a
  ser el puente central entre tasas y montos según `ARCHITECTURE.md`.
- Redondeo bancario en los casos de empate (.5) positivos y negativos.
- Ida y vuelta JSON y SQL, incluida la entrada malformada.
- Descuadre de monedas en toda operación binaria.
- Desbordamiento en `Add`/`Sub`/`Mul` con montos cercanos al límite.

### `decimal` — P2 (ya al 96.6 %, pero es el cimiento)

- Fuzzing de `FromString` y de la aritmética contra `math/big` como oráculo.
- `InexactFloat64` está a 0 % y lo usan casi todas las pruebas del repositorio.
- Rutas de desbordamiento en `Div`, `QuoRem`, `RoundBank`, `divRound`,
  `shr120Round`, `bitAt`, `div3by2`.
- `Pow` con exponentes fraccionarios grandes, bases cercanas a cero,
  exponentes negativos.
- `Log`/`Ln`/`expFpToDec` en los extremos del rango.

### `finance/loans`, `finance/returns`, `finance/simpleinterest` — P2

Cobertura ya alta (90 %, 85 %, 97 %); consolidar:

- `loans`: `SemiAnnually` y `Annually` (0 %), pago anticipado que cancela el
  préstamo antes de tiempo, comparación de refinanciación donde no conviene
  refinanciar, APR con comisiones que superan el principal.
- `returns`: `Must*` sin cobertura, series vacías o de un solo elemento en
  `Mean`/`variance`, rendimiento −100 %, TWR con un flujo intermedio que anula
  la cartera.
- `simpleinterest`: `Present`/`PresentWithFuture` cuando `1 + i·n ≤ 0`.

### `charts` (módulo aparte) — P2

Renderizado con series vacías, con un solo punto, y con valores negativos.

---

## 5. Infraestructura de pruebas

Lo que hay que construir antes o en paralelo, para que escribir cada prueba
nueva cueste poco:

1. **Paquete de ayudantes compartido** (`internal/testutil`): constructores
   `usd(1000)`, comparación de decimales con tolerancia, aserción de que un
   cuadro de amortización cierra, aserción de que una función no entra en
   pánico.
2. **Umbral de cobertura en CI.** El flujo ya calcula `coverage.out`; añadir un
   paso que falle si la cobertura total baja del umbral acordado (empezar en
   el 85 % actual y subirlo por fases) y publicar el desglose por paquete.
3. **`make test-fuzz`** con una duración corta, más una ejecución larga
   programada (nocturna o semanal) en CI.
4. **Corpus de regresión.** Cada bug de la sección 2 aporta su caso a un
   archivo de pruebas de regresión con referencia a este documento, para que no
   vuelva a aparecer.
5. **`Example` en godoc.** Convertir los fragmentos de los comentarios de
   paquete en ejemplos ejecutables, empezando por `bonds`, `gradients`,
   `investment` y `tvm`.

---

## 6. Fases

**Fase 1 — Corregir lo que ya está roto (P0).**
Los pánicos con tasa 0 % (§2.1) y los ceros silenciosos de conversión de tasas
(§2.2). Cada corrección entra junto a la prueba que la fija. Es lo primero
porque son fallos que devuelven cifras equivocadas a quien use la librería hoy.

**Fase 2 — Paquetes con menor cobertura (P1).**
`gradients`, `investment`, `bonds`, `tvm`, `depreciation`. Incluye ejecutar por
primera vez las rutas de bisección y las funciones `Must*` sin cobertura, y
decidir §2.3, §2.4 y §2.5.

**Fase 3 — Invariantes y reciprocidad transversales.**
Las propiedades de la sección 3 aplicadas a todos los paquetes, más las pruebas
de coherencia entre paquetes (`tvm` ↔ `annuities` ↔ `loans`,
`gradients` ↔ `annuities`, `XIRR` ↔ `IRR`).

**Fase 4 — Fuzzing, ejemplos y consolidación (P2).**
`decimal`, `money`, ejemplos de godoc, benchmarks para los paquetes que aún no
los tienen, umbral de cobertura en CI.

---

## 7. Criterios de aceptación

El plan está terminado cuando:

- Ninguna función que devuelva `error` entra en pánico con entradas válidas
  pero degeneradas (tasa 0, un período, monto 0, series vacías).
- Ninguna función devuelve el valor cero con `err == nil` por no haber entrado
  en ninguna rama: todo camino termina en un resultado calculado o en un error
  explícito.
- Todos los tipos enumerados públicos están recorridos por completo en las
  pruebas de su función.
- Ninguna función exportada queda en 0 % de cobertura.
- Las rutas de respaldo de los solucionadores (bisección) se ejecutan al menos
  una vez en las pruebas.
- Cada paquete tiene al menos un `Example` ejecutable.
- Cobertura total ≥ 92 %, sin ningún paquete por debajo del 85 %.

La cobertura va al final de la lista a propósito: es la consecuencia de haber
cubierto los casos de uso, no la meta.
