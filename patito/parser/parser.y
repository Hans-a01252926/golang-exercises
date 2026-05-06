%{
package parser

import (
    "patito/token"
)
}%

%union {
	lit string
}

%token PROGRAMA VARS INICIO FIN
%token ENTERO FLOTANTE NULA
%token SI SINO MIENTRAS HAZ ESCRIBE

%token <lit> ID
%token <lit> CTE_ENT
%token <lit> CTE_FLOT
%token <lit> LETRERO

%token ASIGNA
%token MAS MENOS
%token MULT DIVIDE
%token MAYOR MENOR DIF IGUAL

%token PUNTOCOMA DOSPUNTOS COMA
%token PAR_IZQ PAR_DER
%token LLAVE_IZQ LLAVE_DER

%left MAYOR MENOR DIF IGUAL
%left MAS MENOS
%left MULT DIVIDE
%right UPLUS UMINUS

%start input

%%

input:
	programa
	;

programa:
	PROGRAMA ID PUNTOCOMA vars funcs INICIO cuerpo FIN
	;

vars:
	VARS lista_ids DOSPUNTOS tipo PUNTOCOMA vars
	| /* empty */
	;

lista_ids:
	ID lista_ids_prima
	;

lista_ids_prima:
	COMA ID lista_ids_prima
	| /* empty */
	;

tipo:
	ENTERO
	| FLOTANTE
	;

funcs:
	func funcs
	| /* empty */
	;

func:
	tipo_retorno ID PAR_IZQ params PAR_DER LLAVE_IZQ vars cuerpo LLAVE_DER PUNTOCOMA
	;

tipo_retorno:
	tipo
	| NULA
	;

params:
	ID DOSPUNTOS tipo params_prima
	| /* empty */
	;

params_prima:
	COMA ID DOSPUNTOS tipo params_prima
	| /* empty */
	;

cuerpo:
	LLAVE_IZQ estatutos LLAVE_DER
	;

estatutos:
	estatuto estatutos
	| /* empty */
	;

estatuto:
	asigna
	| condicion
	| ciclo
	| llamada PUNTOCOMA
	| imprime
	;

asigna:
	ID ASIGNA expresion PUNTOCOMA
	;

condicion:
	SI PAR_IZQ expresion PAR_DER cuerpo sino_opc PUNTOCOMA
	;

sino_opc:
	SINO cuerpo
	| /* empty */
	;

ciclo:
	MIENTRAS PAR_IZQ expresion PAR_DER HAZ cuerpo PUNTOCOMA
	;

llamada:
	ID PAR_IZQ argumentos PAR_DER
	;

argumentos:
	expresion argumentos_prima
	| /* empty */
	;

argumentos_prima:
	COMA expresion argumentos_prima
	| /* empty */
	;

imprime:
	ESCRIBE PAR_IZQ imprime_args PAR_DER PUNTOCOMA
	;

imprime_args:
	imprime_val imprime_args_prima
	;

imprime_args_prima:
	COMA imprime_val imprime_args_prima
	| /* empty */
	;

imprime_val:
	expresion
	| LETRERO
	;

expresion:
	exp
	| exp op_rel exp
	;

op_rel:
	MAYOR
	| MENOR
	| DIF
	| IGUAL
	;

exp:
	exp MAS termino
	| exp MENOS termino
	| termino
	;

termino:
	termino MULT factor
	| termino DIVIDE factor
	| factor
	;

factor:
	PAR_IZQ expresion PAR_DER
	| MAS factor %prec UPLUS
	| MENOS factor %prec UMINUS
	| ID
	| cte
	| llamada
	;

cte:
	CTE_ENT
	| CTE_FLOT
	;

%%

func Parse(l yyLexer) int {
	return yyParse(l)
}