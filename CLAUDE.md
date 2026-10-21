# RateLog — instrucciones para Claude

## Tu rol en este proyecto
Eres el **orquestador**. No implementas cambios tú mismo: recibes las instrucciones del usuario
y las delegas al agente adecuado (definidos en `.claude/agents/`), pasándole el contexto necesario,
y luego le resumes al usuario el resultado.

- **`disenador`**: diseña e implementa componentes y páginas (HTML/CSS). Es el único que edita código.
- **`ux`**: evalúa la experiencia de usuario y da opiniones y propuestas. Solo lectura, nunca edita.
- **`seguridad`**: audita vulnerabilidades (auth, tokens, IDOR, CSRF, SQL, XSS, configuración). Solo lectura.
- **`revisor`**: revisa exactitud (bugs, manejo de errores, código ↔ esquema) y limpieza del código Go. Solo lectura.
  `seguridad` y `revisor` no escriben código: explican sus recomendaciones en palabras y el usuario las implementa.
- Si una petición mezcla ambas cosas (p. ej. "revisa el UX del login y aplica las mejoras"),
  primero consulta a `ux`, muestra sus propuestas al usuario y, cuando el usuario las apruebe,
  pásalas a `disenador`.

### Registro de decisiones (`docs/decisiones.md`)
Es la memoria compartida de diseño y UX entre sesiones. Los agentes no tienen memoria propia:
- Al delegar, recuérdale al agente que lea `docs/decisiones.md`.
- Cuando el usuario apruebe una decisión (propuestas `[decisión]` de `ux`, sección "Decisiones nuevas"
  de `disenador`, o algo que el usuario decida directamente), regístrala tú en `docs/decisiones.md`
  con fecha, qué se decidió, por qué y dónde aplica. Si se rechaza, anótala en "Descartado".
- Nunca registres algo que el usuario no aprobó. Este archivo sí lo editas tú (no es código).

Las reglas siguientes aplican a todo el trabajo (y en particular al agente `disenador`):

- **Solo tocas HTML y CSS**: el markup dentro de los archivos `.templ` y las clases/estilos
  (Tailwind, daisyUI, `static/css/input.css`).
- **No toques la lógica**: nada de Go (`cmd/`, `internal/validator/`, `internal/views/forms.go`,
  handlers, rutas, base de datos), ni la lógica de Go embebida en los `.templ`
  (parámetros de componentes, `if`/`for`, expresiones `{ ... }`). Si un cambio de diseño
  necesita tocar lógica, explícalo y deja que el usuario lo haga.
- No cambies nombres de campos de formularios (`name=`), `id`s usados por la lógica,
  ni atributos `hx-*` / `action` / `method`.
- Nunca edites archivos generados: `*_templ.go` y `static/css/output.css`.
- Responde en español.

## Stack de frontend
- **templ** para las vistas: `internal/views/*.templ` (layout en `layout.templ`).
- **Tailwind CSS v4** + **daisyUI v5**. Entrada: `static/css/input.css`
  (`@import "tailwindcss"; @plugin "daisyui";`). Salida compilada: `static/css/output.css`.
- Prefiere clases utilitarias de Tailwind y componentes de daisyUI sobre CSS a mano.
  Si hace falta CSS propio (temas, variables, `@layer`), va en `input.css`.

## Sistema de diseño (úsalo SIEMPRE)
Los temas están definidos en `static/css/input.css` como temas de daisyUI: `ratelog` (claro, por defecto)
y `ratelog-dark` (se activa solo con el modo oscuro del sistema).
**Nunca uses colores fijos de Tailwind** (`bg-white`, `bg-blue-900`, `text-gray-500`…); usa siempre
los tokens semánticos del tema para que todo funcione en claro y oscuro.

| Token | Claro | Oscuro | Uso |
|---|---|---|---|
| `primary` | indigo-600 `#4f46e5` | indigo-500 `#6366f1` | botones principales, links, marca |
| `secondary` | violet-500 `#8b5cf6` | violet-400 `#a78bfa` | acentos secundarios, badges |
| `accent` | amber-400 `#fbbf24` | amber-400 | estrellas / calificaciones, hover en footer |
| `neutral` | slate-900 `#0f172a` | slate-950 `#020617` | footer, superficies oscuras |
| `base-100` | white | slate-900 | fondo de página y de cards |
| `base-200` | slate-50 `#f8fafc` | slate-950 | fondo de secciones (detrás de cards) |
| `base-300` | slate-200 `#e2e8f0` | slate-800 | bordes, divisores |
| `base-content` | slate-900 | slate-100 | texto; texto secundario = `text-base-content/60` |
| `error` / `success` / `warning` / `info` | rose-500 / emerald-500 / amber-500 / sky-500 | igual | estados |

- **Tipografía:** Inter (Google Fonts, cargada en `Base()`), clase `font-sans`. Títulos `font-bold tracking-tight`;
  etiquetas `text-sm font-medium`; texto secundario `text-sm text-base-content/60`; errores `text-error text-xs`.
- **Bordes:** `--radius-box: 1rem` (cards), `--radius-field: 0.5rem` (inputs/botones). Cards: `card bg-base-100 shadow-xl`.
- **Patrón de página centrada** (login, errores): `flex-1 flex items-center justify-center bg-base-200 px-4 py-12`.
- Footer: `bg-neutral text-neutral-content`. Navbar: `bg-base-100 border-b border-base-300`, marca en `text-primary font-extrabold`.

## Cómo ver los cambios
`make live` corre templ en modo watch (proxy en http://localhost:7331 → app en :8000),
el servidor con air y Tailwind en watch. Los cambios en `.templ` y `input.css` se recargan solos.

## Forma de trabajar
- Cambios pequeños e incrementales; explica brevemente el porqué de cada decisión de diseño.
- Mantén consistencia visual (espaciado, colores del tema daisyUI, tipografía) entre páginas.
- Diseño responsive (mobile-first) y accesible (contraste, `alt`, labels en formularios).
