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
  (no solo en hover). Foco: `focus-visible:outline-2 outline-offset-2 outline-primary-content`.
  *Por qué:* contraste garantizado sobre la portada, se descubre en móvil/táctil y funciona con teclado.
  *Dónde:* tarjetas de `MyGamesPage` (botón "Rate and Review") y futuras tarjetas de portada.
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
- **Acciones destructivas en listas** (2026-10-20). Se confirman con un popover (Popover API) con "Cancel" (`btn-ghost`)
  y el botón de confirmar `btn btn-error btn-sm`. Los botones repetidos por elemento (Edit, Delete) llevan el nombre
  del elemento en `<span class="sr-only">` (p. ej. "Delete review of <Name>").

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

## Textos y microcopy

- **Placeholders con contexto** (2026-10-20). Cuando la página trata de un juego concreto, los
  placeholders lo mencionan (p. ej. "What did you like about <Name>? What could be better?").
  Se evitan subtítulos que repiten lo que ya dicen la etiqueta y el título.
- **Mayúsculas** (2026-10-20). Acciones en sentence case ("Log in", "Sign up", "Log out", "Search games");
  nombres de sección en Title Case ("My Games", "My Reviews").
- **Idioma de la UI:** los textos de la interfaz van en inglés, como el resto de la app.

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

## Descartado

Ideas que se evaluaron y se rechazaron (para no volver a proponerlas sin un motivo nuevo).

- **Dropdown con `tabindex` sobre `div`** (2026-10-20): no se cierra al volver a pulsar, ni con Esc, y no
  anuncia su estado. Sustituido por Popover API.
- **"Log out" en rojo** (2026-10-20): cerrar sesión no es destructivo; el rojo se reserva para borrar/quitar.
- **`line-clamp-6` en reviews sin salida** (2026-10-20): ocultaba texto del usuario sin forma de leerlo. Vuelve
  solo cuando exista una página de detalle/edición enlazada desde la tarjeta.
