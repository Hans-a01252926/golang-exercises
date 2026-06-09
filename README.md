# Patito Compiler

Proyecto final de TC3002B.503 _Desarrollo de aplicaciones avanzadas de ciencias computacionales_ para el lenguaje **Patito**, desarrollado en Go.
El proyecto implementa las etapas principales de un compilador académico: análisis léxico, análisis sintáctico, análisis semántico, generación de cuádruplos, direcciones virtuales, archivo objeto y ejecución mediante una Máquina Virtual.

## Descripción

Patito es un mini-lenguaje de programación diseñado para practicar la construcción de compiladores. El compilador toma un archivo `.patito`, valida su sintaxis y semántica, genera código intermedio en forma de cuádruplos, asigna direcciones virtuales y finalmente ejecuta el programa mediante una Máquina Virtual propia.

El lenguaje soporta:

* Variables globales y locales
* Tipos `entero`, `flotante`, `string` y `bool`
* Expresiones aritméticas y relacionales
* Asignaciones
* Impresión con `escribe`
* Condiciones con `si` / `sino`
* Ciclos con `mientras`
* Estatuto `repite`
* Funciones `nula`
* Funciones con retorno
* Parámetros
* Recursión simple
* Algunos casos de recursión doble usando variables locales auxiliares
* Generación de archivo objeto `.obj`
* Ejecución en Máquina Virtual

## Estructura general de un programa Patito

```txt
programa id;

vars globales : tipo;

tipo_retorno funcion(param : tipo) {
    vars locales : tipo;

    {
        estatutos
    }

    return(expresion);
};

inicio {
    estatutos
} fin
```

Ejemplo:

```txt
programa prueba;

vars x, y : entero;
vars texto : string;

entero suma(a : entero, b : entero) {
    {
    }

    return(a + b);
};

inicio {
    x = 2;
    y = suma(x, 3);
    texto = "resultado";

    escribe(texto);
    escribe(y);
} fin
```

## Tipos soportados

| Tipo       | Descripción                                          |
| ---------- | ---------------------------------------------------- |
| `entero`   | Números enteros                                      |
| `flotante` | Números con punto decimal                            |
| `string`   | Cadenas de texto                                     |
| `bool`     | Valores booleanos: `verdadero` y `falso`             |
| `nula`     | Tipo de retorno para funciones que no regresan valor |

## Estructura del proyecto

```txt
patito/
├── lexer/
│   ├── lexer.go
│   └── lexer_test.go
├── parser/
│   ├── parser.y
│   ├── parser.go
│   └── parser_lexer.go
├── token/
│   └── token.go
├── semantic/
│   ├── types.go
│   ├── semantic_cube.go
│   ├── context.go
│   ├── function_directory.go
│   ├── var_table.go
│   └── constants_table.go
├── quadruples/
│   ├── quadruple.go
│   ├── generator.go
│   └── stack.go
├── memoria/
│   └── address_allocator.go
├── maquinavirtual/
│   ├── memory.go
│   ├── execution_memory.go
│   ├── activation_record.go
│   └── maquina_virtual.go
├── objectfile/
│   └── object_file.go
├── pruebas/
├── pruebas_vm/
├── main.go
├── go.mod
└── README.md
```

> Los nombres exactos de carpetas pueden variar según la versión local del proyecto, pero la organización general separa lexer, parser, semántica, cuádruplos, memoria, archivo objeto y Máquina Virtual.

## Componentes principales

### Lexer

El lexer se encarga de leer el código fuente y convertirlo en tokens.
Reconoce palabras reservadas, identificadores, constantes, operadores y signos de puntuación.

Ejemplos de tokens:

```txt
PROGRAMA
VARS
ID
CTE_ENT
CTE_FLOT
LETRERO
SI
SINO
MIENTRAS
ESCRIBE
RETURN
```

### Parser

El parser fue construido con `goyacc`.

El archivo principal es:

```txt
parser/parser.y
```

A partir de este archivo se genera:

```txt
parser/parser.go
```

El parser valida que el programa siga la gramática de Patito y ejecuta acciones semánticas durante el reconocimiento de reglas.

### Parser Lexer Adapter

El archivo `parser_lexer.go` funciona como puente entre el lexer manual y el parser generado por `goyacc`.

El lexer produce tokens del paquete `token`, mientras que `goyacc` espera sus propios tokens internos. El adaptador convierte los tokens del lexer al formato que entiende el parser.

También guarda el contexto semántico y el generador de cuádruplos para que las acciones de `parser.y` puedan usarlos.

### Análisis semántico

El análisis semántico valida que el programa tenga sentido más allá de la sintaxis.

Incluye:

* Cubo semántico
* Tabla de variables
* Directorio de funciones
* Tabla de constantes
* Validación de tipos
* Validación de variables declaradas
* Validación de funciones declaradas
* Validación de parámetros
* Validación de retornos

### Cubo semántico

El cubo semántico valida operaciones entre tipos.

Conceptualmente responde:

```txt
tipo izquierdo + operador + tipo derecho = tipo resultado
```

Ejemplo:

```txt
entero + entero     → entero
entero + flotante   → flotante
flotante > entero   → bool
string == string    → bool
bool == bool        → bool
entero = flotante   → error
```

### Cuádruplos

Los cuádruplos son la representación intermedia del programa.

Formato:

```txt
(operador, left, right, result)
```

Ejemplo:

```txt
(+, 1000, 1001, 9000)
```

Esto significa:

```txt
leer 1000
leer 1001
sumar
guardar resultado en 9000
```

Operadores principales:

| Categoría    | Operadores                                   |
| ------------ | -------------------------------------------- |
| Aritméticos  | `+`, `-`, `*`, `/`                           |
| Relacionales | `<`, `>`, `==`, `!=`                         |
| Control      | `=`, `print`, `GOTO`, `GOTOF`                |
| Funciones    | `ERA`, `PARAM`, `GOSUB`, `RETURN`, `ENDFunc` |

### Direcciones virtuales

Las direcciones virtuales permiten que los cuádruplos trabajen con direcciones de memoria en lugar de nombres de variables.

Distribución final:

| Segmento  | Entero        | Flotante      | String        | Bool          |
| --------- | ------------- | ------------- | ------------- | ------------- |
| Global    | 1000 - 1999   | 2000 - 2999   | 3000 - 3999   | 4000 - 4999   |
| Local     | 5000 - 5999   | 6000 - 6999   | 7000 - 7999   | 8000 - 8999   |
| Temporal  | 9000 - 9999   | 10000 - 10999 | 12000 - 12999 | 11000 - 11999 |
| Constante | 13000 - 13999 | 14000 - 14999 | 15000 - 15999 | 16000 - 16999 |

Ejemplo:

```txt
x = 5;
```

Puede generar:

```txt
(=, 13000, _, 1000)
```

Donde:

```txt
13000 → constante 5
1000  → variable x
```

### Archivo objeto

El compilador genera un archivo `.obj` en formato JSON.

Este archivo contiene:

* Lista de cuádruplos
* Tabla de constantes

Flujo:

```txt
archivo.patito
→ scanner
→ parser
→ análisis semántico
→ cuádruplos
→ archivo.obj
→ Máquina Virtual
→ ejecución
```

### Máquina Virtual

La Máquina Virtual carga el archivo objeto y ejecuta los cuádruplos.

La memoria de ejecución se divide en:

```txt
ExecutionMemory
├── GlobalMemory
├── ConstMemory
├── TempMemory
├── CallStack
└── PendingERA
```

Cada llamada a función usa un `ActivationRecord`, que contiene:

```txt
ActivationRecord
├── FunctionName
├── LocalMemory
├── TempMemory
└── ReturnIP
```

Esto permite ejecutar funciones con parámetros, funciones con retorno y recursión.

## Comandos útiles

### Generar el parser con goyacc

```bash
goyacc -o parser/parser.go -v y.output parser/parser.y
```

### Ejecutar un programa Patito

```bash
go run . pruebas_vm/archivo.patito
```

Ejemplo:

```bash
go run . pruebas_vm/factorial_recursivo.patito
```

### Ejecutar pruebas del lexer

```bash
go test ./lexer
```

### Ejecutar todas las pruebas de Go

```bash
go test ./...
```

## Ejemplo de programa

```txt
programa pruebaFactorial;

vars resultado : entero;

entero factorial(n : entero) {
    vars r : entero;

    {
        si (n < 1) {
            r = 1;
        } sino {
            r = n * factorial(n - 1);
        };
    }

    return(r);
};

inicio {
    resultado = factorial(5);
    escribe(resultado);
} fin
```

Salida esperada:

```txt
120
```

## Funciones

Patito maneja funciones con los cuádruplos:

```txt
ERA
PARAM
GOSUB
RETURN
ENDFunc
```

### ERA

Prepara el registro de activación de la función.

### PARAM

Copia un argumento al parámetro formal.

### GOSUB

Salta al inicio de la función.

### RETURN

Guarda el valor de retorno.

### ENDFunc

Termina la función y regresa al punto de llamada.

## Limitaciones conocidas

Patito todavía tiene algunas limitaciones:

* No maneja arreglos.
* No tiene lectura de input desde consola.
* No tiene operadores lógicos `&&`, `||` o `!`.
* No tiene operadores `<=` ni `>=`.
* No maneja múltiples `return` dentro de `if/else`.
* El `return` se maneja al final de la función.
* No tiene funciones estándar como `sqrt`, `pow` o `len`.
* No tiene comentarios.

## Notas sobre recursión

Patito soporta recursión simple, como factorial.

Para recursión doble, como Fibonacci, se recomienda usar variables locales auxiliares para guardar cada llamada recursiva antes de combinar resultados.

Ejemplo:

```txt
entero fibonacci(n : entero) {
    vars r, a, b : entero;

    {
        si (n < 2) {
            r = n;
        } sino {
            a = fibonacci(n - 1);
            b = fibonacci(n - 2);
            r = a + b;
        };
    }

    return(r);
};
```

## Autor

Hans Gerhard Moreno

## Créditos

Proyecto desarrollado como parte de la materia de compiladores.

Se utilizó Go, `goyacc` y apoyo de herramientas de IA para organización, depuración y documentación del proyecto.

Además, se tomó como referencia el libro **_Writing an Interpreter in Go_** de **Thorsten Ball**, especialmente para entender la implementación de un lexer manual, la definición de tokens y la estructura inicial de un intérprete en Go.
