# Decisiones de diseño y UX — RateLog

Registro de las decisiones de diseño y experiencia de usuario que ya se aprobaron.
Las bases del sistema de diseño (tokens, tipografía, bordes, patrones de página) están en `CLAUDE.md`;
aquí van las decisiones que se toman encima de esas bases mientras se construye la app.

**Reglas**
- Lo leen siempre los agentes `disenador` y `ux` antes de trabajar, y lo respetan.
- Solo lo actualiza el orquestador, y solo con decisiones que el usuario aprobó.
- Si una decisión cambia, se edita su entrada (no se duplica) y se anota la fecha del cambio.
- Cada entrada: qué se decidió, por qué y dónde aplica.

## Componentes y patrones

- **Acciones en tarjetas con portada** (2026-10-20). El botón va dentro del degradado inferior
  `from-neutral`, bajo el título: `btn btn-primary btn-xs sm:btn-sm w-full`. Siempre visible
  (no solo en hover). Foco: `focus-visible:outline-2 outline-offset-2 outline-neutral-content`
  (cambiado el 2026-10-21: antes `outline-primary-content`, que dejó de verse sobre `neutral` en oscuro al pasar
  `primary-content` a slate-900).
  *Por qué:* contraste garantizado sobre la portada, se descubre en móvil/táctil y funciona con teclado.
  *Dónde:* tarjetas de `MyGamesPage` (botón "Rate and review") y futuras tarjetas de portada.
- **Calificaciones con estrellas** (2026-10-20). Componente `StarScore()`: `rating` de daisyUI con
  radios `name="rating"` (1–5), `mask mask-star-2 bg-accent`, `aria-label` por estrella, dentro de
  `fieldset` + `legend`. Foco visible en el contenedor (`has-focus-visible:outline-primary`, `rounded-field`).
  *Por qué:* `accent` es el token de calificaciones; los radios nativos dan teclado gratis; la máscara
  recorta el outline del input. *Dónde:* toda calificación de la app.
- **Textareas** (2026-10-20). `textarea w-full min-h-36 leading-relaxed`, `label` `text-sm font-medium`
  encima, separación `gap-1.5` como los inputs. *Dónde:* todo campo de texto largo.
- **Cabecera de juego** (2026-10-20). En páginas sobre un juego concreto, la portada va como cabecera
  de la card: `figure` `aspect-[2/1] bg-base-300`, imagen `object-cover`, degradado inferior
  `from-base-100 to-transparent`, `overflow-hidden` en la card, `alt="Cover art for <Name>"`, sin `loading="lazy"`.
  El `h1` es el nombre del juego (`text-balance`); la acción va encima como etiqueta
  `text-xs font-semibold uppercase tracking-wider text-primary`.
  *Por qué:* el usuario reconoce el juego al instante; 2:1 deja el formulario visible en móvil sin scroll.
  *Dónde:* `ReviewPage` y futuras páginas de un juego.
- **Datos del juego** (2026-10-20). Línea secundaria `text-sm text-base-content/60` separada con `·`
  (`aria-hidden`). Fecha solo como año en `<time datetime="YYYY-MM-DD">`; números con `tabular-nums`.
  *Por qué:* da contexto sin competir con el título; el año distingue remakes.

- **Errores de validación** (2026-10-20). El mensaje va debajo del campo, dentro del mismo contenedor
  `flex flex-col gap-1.5` que el label y el campo: `<p class="text-error text-xs">`. En componentes
  agrupados (estrellas u otro `fieldset`), se envuelve en `div flex flex-col items-center gap-1.5` con el
  error justo después, con la misma alineación que el componente.
  *Por qué:* consistencia con `LoginForm`/`SignupForm`; el error queda pegado al campo que lo causa.
  *Dónde:* todos los formularios.

- **Navbar** (2026-10-20). `<header class="bg-base-100 border-b border-base-300">` + `<nav aria-label="Main" class="navbar">`.
  Marca "Rate Log" a la izquierda (`text-primary font-extrabold`) enlazando a `/`. Escritorio (`lg:`):
  `menu menu-horizontal` con icono + texto (Search games, My Games, My Reviews). Por debajo de `lg`: una
  sola hamburguesa a la derecha que abre el menú de cuenta; dentro, primero la navegación (`lg:hidden`,
  `py-2.5`), divisor, "Signed in as {nombre}" y "Log out". Nunca iconos sueltos sin texto.
  *Dónde:* `NavBar` en `layout.templ`. `BaseLogin()` lleva una cabecera mínima solo con la marca.
- **Desplegables** (2026-10-20). Siempre con Popover API de daisyUI v5 (`<button popovertarget=… style="anchor-name:…">`
  + `<ul popover class="dropdown dropdown-end menu …" style="position-anchor:…">`) o `<details>`. Nunca `tabindex`
  sobre un `div`, ni `role="menu"`/`menuitem`. Contenido: `menu w-64 rounded-box border border-base-300 bg-base-100 p-2 shadow-xl`.
  *Por qué:* se cierran al volver a pulsar, con Esc y al tocar fuera, y exponen `aria-expanded` sin JS.
  *Ojo:* sin soporte de anchor positioning el popover sale centrado; probar en Firefox y Safari.
- **Divisor en menús** (2026-10-20). `<li aria-hidden="true" class="my-1 h-px bg-base-300">` en vez de bordes sobre un ítem,
  para poder ocultarlo por breakpoint sin afectar a los vecinos.
- **Acciones POST dentro de menús** (2026-10-20). `<form id="…" method="post" action="…" class="hidden"></form>` fuera
  de la lista y `<button type="submit" form="…">` directamente en el `li`, para que herede el estilo de `menu`.
- **Iconos por sección** (2026-10-20). Lupa = Search games, pila/colección = My Games, bocadillo de texto = My Reviews.
  Siempre `aria-hidden="true"` y acompañados de texto. La estrella se reserva para puntuar.
- **Color `error`** (2026-10-20). `text-error`/`btn-error` solo para acciones destructivas o irreversibles (borrar,
  quitar). Cerrar sesión no lo es: va con color de texto normal, al final y separado por el divisor.

- **Notas de solo lectura** (2026-10-20). `div.rating` (`rating-sm` en listas) con `role="img"` y
  `aria-label="Rated N out of 5"`; cinco `div.mask.mask-star-2.bg-accent` con `aria-hidden` y `aria-current?={ rating == N }`;
  al lado el texto visible "N/5" (`text-sm font-medium tabular-nums`, `aria-hidden`). Nunca radios para mostrar una nota.
  La regla de `input.css` para `aria-current` sin valor es provisional (sobra si se usa `strconv.FormatBool`).
- **Tarjetas de review** (2026-10-20). Rejilla de 1 columna, 2 desde `lg`, tarjetas de igual alto (`h-full`),
  `ul role="list"`. Portada lateral desde `sm` (`sm:card-side`, `w-40`/`md:w-48`) sin estirarse con reviews largas;
  en móvil, portada arriba en `aspect-[3/1]` (el texto conserva todo el ancho); desde `sm`,
  `sm:aspect-auto sm:w-40 md:w-48 sm:min-h-40 sm:max-h-56`. La miniatura lateral en móvil queda para
  tarjetas con poco texto (como `GameCard`). Orden interno: título → nota → texto → datos.
- **Texto largo del usuario** (2026-10-20). `whitespace-pre-line break-words leading-relaxed text-base-content/80`.
  Nunca truncar contenido del usuario (`line-clamp`) sin un camino para verlo completo (enlace a detalle/edición).
- **Portadas en listas** (2026-10-20). `loading="lazy"` y `alt=""` cuando el nombre del juego está junto a la imagen
  (decorativa). La "Cabecera de juego" de una página individual sigue con `alt` descriptivo y sin lazy.
- **Fechas de actividad** (2026-10-20). Formato `Jan 2, 2006` en `<time datetime="YYYY-MM-DD">` dentro de la línea de
  datos (`tabular-nums`); el año de lanzamiento sigue siendo solo el año.
- **Estados vacíos** (2026-10-20). Toda lista tiene estado vacío: card `bg-base-100 shadow-xl` centrada con el icono
  de la sección (`text-base-content/40`, `aria-hidden`), `h2` "No … yet", una línea `text-sm text-base-content/60`
  y un único CTA `btn btn-primary` que lleva al paso anterior del flujo. No se renderiza el `<ul>` vacío.
  Medidas (2026-10-22): card `mx-auto w-full max-w-lg`, `card-body items-center gap-3 py-12 text-center`, icono `h-12 w-12`
  con `stroke-width="1.5"`, `h2 text-xl font-bold tracking-tight`, texto con `text-balance`; va dentro de un contenedor
  con el mismo margen superior que tendría la lista.
- **Acciones destructivas en listas** (2026-10-20). Se confirman con un popover (Popover API) con "Cancel" (`btn-ghost`)
  y el botón de confirmar `btn btn-error btn-sm`. Los botones repetidos por elemento (Edit, Delete) llevan el nombre
  del elemento en `<span class="sr-only">` (p. ej. "Delete review of <Name>").
  *Excepción (2026-10-22):* en tarjetas con portada sin ids por elemento se usa "Confirmación dentro de la tarjeta" (`<details>`).
- **Error de formulario** (2026-10-21). Los errores que no son de un campo concreto (credenciales inválidas,
  fallos generales del envío) se muestran una sola vez encima del primer campo, nunca repetidos ni bajo un campo:
  `<div id="form-error" class="alert alert-error alert-soft text-sm">` con icono `aria-hidden`, sin `role="alert"`;
  el `<form>` lleva `aria-describedby="form-error"` cuando existe. En la lógica es una clave propia (`"form"`).
  *Por qué:* un error de credenciales bajo Email hace creer que el email está mal.
  *Dónde:* login y cualquier formulario con errores generales.
  *Posición (2026-10-27):* en cualquier formulario dentro de una card va en el mismo sitio que los avisos flash:
  dentro del `card-body`, tras el bloque de título (y cabecera de juego si la hay) y antes del primer campo. Si
  coinciden un aviso flash y un error de formulario, el error va primero.
- **Iconos de los avisos** (2026-10-27). Un icono de trazo por tipo de `alert` (heroicons, `stroke-width="2"`,
  `h-5 w-5 shrink-0`, `aria-hidden`, color heredado con `currentColor`): `error` = exclamation-circle,
  `success` = check-circle, `info` = information-circle, `warning` = exclamation-triangle.
  *Por qué:* el mismo icono en todos los `#form-error` y avisos evita que cada página invente el suyo.
  *Dónde:* todo `#form-error` (reviews, login, signup) y los avisos flash.
- **Avisos flash** (2026-10-21). Tras una redirección, aviso dentro de `card-body`, justo después del bloque
  `h1` + subtítulo y antes del primer campo, a todo el ancho: `<div class="alert alert-soft alert-<tipo> text-sm">`
  con icono `aria-hidden` y el texto en un `<span>`. Sin `role` (igual que "Anuncio de errores"), no desaparece solo
  ni lleva botón de cerrar. La lógica usa una cookie de un solo uso con un código de lista cerrada; el texto lo elige
  el servidor. Textos: registro → `success` "Account created. Log in to start rating games."; logout → `info`
  "You've been logged out."; página protegida sin sesión → `info` "Please log in to continue."; sesión caducada o
  token inválido → `warning` "Your session has expired. Log in again to continue."
  *Por qué:* el usuario sabe qué pasó tras cada redirección; sin temporizador no obliga a leer deprisa (WCAG 2.2.1).
- **Tipos de alerta** (2026-10-21). `success` = algo completado, `info` = estado neutro, `warning` = algo que el
  usuario no provocó, `error` = solo errores de formulario.
- **Repoblar formularios con errores** (2026-10-21). Al volver a mostrar un formulario con errores, cada campo
  recupera lo enviado: inputs con `value={ values.X }`, textareas con `{ values.X }` como contenido, radios/checkboxes
  con `checked?={ … }`. Los componentes reutilizables (p. ej. `StarScore`) reciben su valor actual como parámetro
  simple, no el struct del formulario. Las contraseñas nunca se repueblan.
- **Validación del servidor** (2026-10-21). Los `<form>` llevan `novalidate`: todos los errores vienen del servidor
  con el mismo aspecto, sin los globos nativos del navegador. *Excepción:* el buscador (`games.templ`) mantiene la
  validación nativa (`required`) porque es una búsqueda sin errores del servidor. No aplica a formularios ocultos de solo acción.
- **Estructura de campos con icono** (2026-10-21). Un solo `<label for="x" class="flex flex-col gap-1.5">` con
  `<span class="text-sm font-medium">` (texto) y `<span class="input w-full">` (icono `aria-hidden` + input). Error y
  texto de ayuda van fuera del label, dentro del `div flex flex-col gap-1.5`. Nunca dos labels para el mismo input.
  *Por qué:* el clic en texto, icono o borde enfoca el input y el nombre accesible es solo el texto.
  *Dónde:* todos los formularios, incluido el buscador.
- **Autocompletado** (2026-10-21). Todo input de usuario lleva `autocomplete` explícito (`email`, `current-password`,
  `new-password`, `nickname`…); el primer campo de un formulario de acceso lleva `autofocus`.

- **Avisos flash en páginas de sección** (2026-10-22). En páginas de sección (Search games, My Games…) el aviso flash
  va justo debajo del `h1` + subtítulo y encima del contenido, con el mismo ancho, con el estilo de "Avisos flash"
  (no dentro de un `card-body`). Textos al añadir un juego: `success` "<Name> added to My Games." /
  `info` "<Name> is already in My Games." (la lógica decide cuál).
- **Excepción "sin resultados" a "Estados vacíos"** (2026-10-22). Una búsqueda sin resultados muestra la card de estado
  vacío con la lupa y sin `btn-primary`: como mucho un enlace "Clear search", porque el paso anterior del flujo es el
  propio buscador, que ya está en la página. Se distinguen tres estados: sin buscar, sin resultados, con resultados.
- **Estado "In My Games"** (2026-10-22). Un elemento que ya está en una colección del usuario muestra su estado en lugar
  de la acción de añadir: enlace `btn btn-ghost` del mismo tamaño que la acción a la que sustituye (en `GameCard`,
  `btn-sm sm:btn-md`; cambiado el 2026-10-22, antes `btn-sm` fijo) con icono de check `aria-hidden` y texto "In My Games" que lleva
  a `/my_games`, con el nombre del elemento en `sr-only`. Nunca un botón deshabilitado. Requiere que la lógica pase el
  conjunto de ids que el usuario ya tiene.
- **Botones repetidos en listas** (2026-10-22). Amplía "Acciones destructivas en listas": todo botón o enlace que se
  repite por elemento (Add, Rate and review, Edit, Delete…) lleva el nombre del elemento en `<span class="sr-only">`.
  Nunca un `aria-label` que sustituya al texto visible.
- **Feedback en envíos lentos** (2026-10-22). Solo en acciones que dependen de servicios externos (buscar y añadir
  juegos, RAWG): al enviar, el botón muestra `loading loading-spinner loading-xs`, `aria-busy="true"` y el texto
  "Searching…"/"Adding…". Requiere un mínimo de JS o HTMX (lógica).

- **Hover en tarjetas no clicables** (2026-10-22). Las tarjetas cuya acción es un botón interno no llevan elevación ni
  zoom de imagen al pasar el ratón; solo reacciona el botón. Solo una tarjeta que sea entera un enlace puede elevarse.
  *Por qué:* el movimiento invitaba a pulsar zonas que no hacen nada. *Dónde:* `GameCard`, tiles de `MyGamesPage`.
- **Estado "sin buscar"** (2026-10-22). No lleva card: basta con `h1`, subtítulo y buscador. Completa los tres estados
  del buscador (sin buscar / sin resultados / con resultados).
- **Acción secundaria en tarjetas con portada** (2026-10-22). Va en la misma fila del degradado que la principal: la
  principal `flex-1` con el texto en `truncate`; la secundaria `btn btn-ghost btn-square btn-xs sm:btn-sm` solo con
  icono, `border-neutral-content/30` y nombre en `sr-only`. Si es destructiva, rojo solo en hover (`hover:text-error`).
  Nunca con el mismo peso visual que la principal.
- **Confirmación dentro de la tarjeta** (2026-10-22). Si no hay ids por elemento, `<details name="…">` (uno abierto a la
  vez). Al abrirse, un panel `absolute inset-0 bg-neutral/95` cubre la tarjeta: pregunta `text-sm font-bold`,
  consecuencia `text-xs text-neutral-content/80`, confirmar `btn-error` a todo el ancho, y el `summary` convertido en
  "Cancel" `btn-ghost` debajo. Lo que queda tapado se oculta con `group-has-[details[open]]/…:invisible`.
  *Limitación:* no se cierra con Esc ni al tocar fuera. *Dónde:* quitar juego en `MyGamesPage`.
- **Icono de quitar/borrar** (2026-10-22). Papelera de trazo (`stroke-width="2"`, `aria-hidden`), reservada para acciones
  destructivas.
- **Confirmar antes en vez de deshacer** (2026-10-22). Las acciones destructivas se confirman antes; no hay "Deshacer"
  (exigiría borrado lógico). Se podrá revisar si algún día se añade borrado lógico.
- **Edición de review** (2026-10-22). Si ya existe una review del juego, la página dice "Edit your review" en la
  etiqueta superior y el botón "Save changes"; en My Games la tarjeta de un juego reseñado muestra "Edit review" en lugar
  de "Rate and review" (con el nombre en `sr-only`).

## Páginas y flujos

- **Inicio canónico** (2026-10-20). `/` es el inicio de la app (redirige a `/search_game`, que manda a `/login` sin
  sesión). Todo enlace a "inicio" (marca, "Back to home") usa `href="/"`.

- **Formularios largos** (2026-10-20). Página centrada con `card w-full max-w-lg bg-base-100 shadow-xl`
  (más ancha que el `max-w-sm` de login/signup). Acciones al final: principal `btn btn-primary` y
  secundaria "Cancel" `btn btn-ghost` que vuelve a la página anterior; apiladas en móvil (principal
  arriba) y alineadas a la derecha desde `sm:`. *Por qué:* espacio para escribir y una salida clara.
  *Dónde:* `ReviewPage` y futuros formularios de más de dos campos.
- **Orden de listas personales** (2026-10-20). De la más reciente a la más antigua (lógica: `ORDER BY … DESC`).
- **Títulos de sección** (2026-10-20). El `h1` de cada sección coincide literalmente con su etiqueta en la navbar
  ("Search games", "My Games", "My Reviews").
- **Título de pestaña por página** (2026-10-21). `<title>` con el formato "<Página> · Rate Log"
  (p. ej. "Log in · Rate Log", "Sign up · Rate Log"). `Base`/`BaseLogin` reciben el título como parámetro (lógica).
- **Registro con inicio de sesión automático** (2026-10-21). Al crear la cuenta se inicia sesión y se redirige a `/`
  (o al `next` válido), en lugar de mandar al login.
- **Volver a la página de destino (`next`)** (2026-10-21). El middleware añade `?next=<RequestURI>` solo en peticiones
  GET. Viaja en `<input type="hidden" name="next">` en el login y se propaga al enlace "Sign up". Válido solo si es una
  ruta local (empieza por una única `/`, sin `\`, sin esquema ni host, sin caracteres de control) y no es `/login`,
  `/signup` ni `/logout`. Si no es válido se ignora en silencio y se redirige a `/` ("Inicio canónico").
- **Códigos de estado de formularios** (2026-10-21). Un formulario que se vuelve a mostrar con errores responde 422.
- **Fallo del servidor al enviar un formulario** (2026-10-21). Se vuelve a mostrar ese formulario con sus datos y un
  error de formulario general (con 500), en vez de la página de error genérica.
- **Rate limit sin fricción** (2026-10-21). Un usuario que navega con normalidad nunca ve un 429: `/static/` queda fuera
  del limitador y el límite estricto se reserva para los envíos de autenticación (POST `/login`, `/signup`).

- **Una sola review por juego, editable** (2026-10-22). Volver a abrir la review de un juego ya reseñado abre la edición
  precargada con su texto y nota; nunca un formulario vacío que sobrescriba la anterior.
- **Recurso ajeno o inexistente = 404** (2026-10-22). Un id no numérico o ≤ 0, un juego inexistente o un juego fuera de
  la biblioteca del usuario responden 404 (nunca 500), en GET y en POST.
- **Aviso al guardar una review** (2026-10-22). Tras guardar se redirige a My Reviews con `success` "Review saved." /
  "Review updated.", según "Avisos flash en páginas de sección".
- **Servicio externo no disponible** (2026-10-22). Si falla RAWG (caído, cuota, timeout) se vuelve a mostrar la misma página
  conservando lo que escribió el usuario, con un aviso `alert alert-soft alert-warning` en el sitio del flash:
  "Game search is temporarily unavailable. Try again in a few minutes." (503). Al añadir, el mismo aviso vía flash.
  Nunca la página genérica de 500.
- **Volver a la lista con su estado** (2026-10-22). Tras una acción en una lista se vuelve a la misma lista con su estado
  (búsqueda, filtros): el formulario de añadir lleva el `title` actual y se redirige a `/search_game?title=…`.
- **Quitar un juego borra su review** (2026-10-22). Decisión del usuario: al quitar un juego de la biblioteca se borra
  también su review (FK en cascada `reviews` → `users_games`). La confirmación lo avisa. Se puede revisar más adelante.

## Textos y microcopy

- **Placeholders con contexto** (2026-10-20). Cuando la página trata de un juego concreto, los
  placeholders lo mencionan (p. ej. "What did you like about <Name>? What could be better?").
  Se evitan subtítulos que repiten lo que ya dicen la etiqueta y el título.
- **Mayúsculas** (2026-10-20). Acciones en sentence case ("Log in", "Sign up", "Log out", "Search games");
  nombres de sección en Title Case ("My Games", "My Reviews").
- **Idioma de la UI:** los textos de la interfaz van en inglés, como el resto de la app.
- **Mensajes de validación** (2026-10-21). Sin términos técnicos ("characters", nunca "bytes") y en sentence case
  (p. ej. "Must be at least 8 characters", "Passwords don't match", "Choose a rating"). Las reglas se muestran
  antes del error como texto de ayuda `text-xs text-base-content/60` bajo el campo (`id="<campo>-hint"`, enlazado
  con `aria-describedby`); cuando hay error, el error ocupa su lugar.

- **Normalizar la búsqueda** (2026-10-22). Antes de validar y buscar, el texto se recorta, se unen los espacios
  repetidos, se quitan acentos y la puntuación cuenta como espacio; la coincidencia exige cada palabra por separado,
  en cualquier orden. Si queda vacío equivale a "sin buscar". *Dónde:* lógica del buscador (`client.go`/handler).

- **"Remove" frente a "Delete"** (2026-10-22). "Remove" para sacar algo de una colección del usuario (el juego sigue en el
  catálogo); "Delete" solo para contenido que el usuario creó y que se pierde (reviews). Nombre accesible completo:
  "Remove <Name> from My Games".
- **Confirmación con consecuencias** (2026-10-22). Si al quitar algo se borran datos que el usuario creó, la confirmación
  lo dice en una línea (p. ej. "Your rating and review of this game will be deleted too.").
- **Avisos de la biblioteca** (2026-10-22). Quitar: `success` "<Name> removed from My Games."; si no estaba: `info`
  "<Name> wasn't in My Games." (o sin aviso si no se conoce el nombre).

## Accesibilidad y responsive

- **Ids de error** (2026-10-20). Cada mensaje de error lleva `id="<campo>-error"` (p. ej. `body-error`,
  `rating-error`) para enlazarlo con el campo vía `aria-describedby`.
- **Anuncio de errores** (2026-10-20). Sin `role="alert"` en errores que llegan con recarga completa de la
  página; se usa `aria-invalid` + `aria-describedby` en el campo (la condición la pone la lógica de Go).
  *Por qué:* un `alert` al cargar la página se anuncia de forma desigual según el lector de pantalla.
- **Foco consistente** (2026-10-20). Todo elemento de navegación y de menú (marca, enlaces, disparador, ítems)
  usa `focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary`.
- **Disparadores solo-icono** (2026-10-20). Nombre accesible con `<span class="sr-only">` siempre presente, que
  incluye el texto visible cuando lo hay (p. ej. "<nombre> Account menu"); sin verbos de estado ("Open…").
- **Contraste de `primary` en oscuro** (2026-10-21). En `ratelog-dark`, `primary` pasa a indigo-400 `#818cf8`
  (~6:1 sobre slate-900) para que los links de texto lleguen a AA. *Dónde:* `static/css/input.css`; sustituye el
  valor indigo-500 de la tabla de tokens de `CLAUDE.md`. `primary-content` en oscuro pasa a slate-900 `#0f172a`
  (~6:1 sobre indigo-400; el blanco daba ~2,98:1), igual que `secondary-content`; en claro sigue blanco.

- **Botones que pasan a solo-icono en móvil** (2026-10-22). El texto va en `sr-only sm:not-sr-only` (nunca `hidden`) y el
  nombre del elemento en otro `<span class="sr-only">` aparte, precedido de ": " o integrado en la frase.

## Descartado

Ideas que se evaluaron y se rechazaron (para no volver a proponerlas sin un motivo nuevo).

- **Dropdown con `tabindex` sobre `div`** (2026-10-20): no se cierra al volver a pulsar, ni con Esc, y no
  anuncia su estado. Sustituido por Popover API.
- **"Log out" en rojo** (2026-10-20): cerrar sesión no es destructivo; el rojo se reserva para borrar/quitar.
- **`line-clamp-6` en reviews sin salida** (2026-10-20): ocultaba texto del usuario sin forma de leerlo. Vuelve
  solo cuando exista una página de detalle/edición enlazada desde la tarjeta.
- **Textos humanos en las páginas de error** (2026-10-21): se propuso sustituir `http.StatusText` ("Not Found",
  "Too Many Requests", "Unprocessable Entity"…) por mensajes como "Page not found" o "Too many requests. Wait a
  moment and try again.". Se deja para más adelante: de momento las páginas de error usan el texto estándar del código.
