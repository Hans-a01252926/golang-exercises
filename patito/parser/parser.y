%{
package parser

import (
	"patito/semantic"
)
%}

%union {
	lit string
	typ semantic.Type
	ids []string
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

%type <typ> tipo tipo_retorno expresion exp termino factor cte
%type <ids> lista_ids lista_ids_prima

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
	VARS lista_ids DOSPUNTOS tipo PUNTOCOMA
	{
		yylex.(*PatitoLexer).Sem.AddVars($2, $4)
	}
	vars
	| /* empty */
	;

lista_ids:
	ID lista_ids_prima
	{
		$$ = append([]string{$1}, $2...)
	}
	;

lista_ids_prima:
	COMA ID lista_ids_prima
	{
		$$ = append([]string{$2}, $3...)
	}
	| /* empty */
	{
		$$ = []string{}
	}
	;
tipo:
	ENTERO
	{
		$$ = semantic.TypeEntero
	}
	| FLOTANTE
	{
		$$ = semantic.TypeFlotante
	}
	;
	
funcs:
	func funcs
	| /* empty */
	;

func:
	tipo_retorno ID
	{
		yylex.(*PatitoLexer).Sem.StartFunction($2, $1)
	}
	PAR_IZQ params PAR_DER LLAVE_IZQ vars cuerpo LLAVE_DER PUNTOCOMA
	{
		yylex.(*PatitoLexer).Sem.EndFunction()
	}
	;

tipo_retorno:
	tipo
	{
		$$ = $1
	}
	| NULA
	{
		$$ = semantic.TypeNula
	}
	;

params:
	ID DOSPUNTOS tipo
	{
		yylex.(*PatitoLexer).Sem.AddVar($1, $3)
	}
	params_prima
	| /* empty */
	;

params_prima:
	COMA ID DOSPUNTOS tipo
	{
		yylex.(*PatitoLexer).Sem.AddVar($2, $4)
	}
	params_prima
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
	{
		yylex.(*PatitoLexer).Sem.CheckAssignment($1, $3)
	}
	;

condicion:
	SI PAR_IZQ expresion PAR_DER
	{
		yylex.(*PatitoLexer).Sem.CheckCondition($3)
	}
	cuerpo sino_opc PUNTOCOMA
	;

sino_opc:
	SINO cuerpo
	| /* empty */
	;

ciclo:
	MIENTRAS PAR_IZQ expresion PAR_DER
	{
		yylex.(*PatitoLexer).Sem.CheckCondition($3)
	}
	HAZ cuerpo PUNTOCOMA
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
	{
		$$ = $1
	}
	| exp MAYOR exp
	{
		$$ = yylex.(*PatitoLexer).Sem.CheckOperation($1, semantic.OpMayor, $3)
	}
	| exp MENOR exp
	{
		$$ = yylex.(*PatitoLexer).Sem.CheckOperation($1, semantic.OpMenor, $3)
	}
	| exp DIF exp
	{
		$$ = yylex.(*PatitoLexer).Sem.CheckOperation($1, semantic.OpDif, $3)
	}
	| exp IGUAL exp
	{
		$$ = yylex.(*PatitoLexer).Sem.CheckOperation($1, semantic.OpIgual, $3)
	}
	;

exp:
	exp MAS termino
	{
		$$ = yylex.(*PatitoLexer).Sem.CheckOperation($1, semantic.OpSuma, $3)
	}
	| exp MENOS termino
	{
		$$ = yylex.(*PatitoLexer).Sem.CheckOperation($1, semantic.OpResta, $3)
	}
	| termino
	{
		$$ = $1
	}
	;


termino:
	termino MULT factor
	{
		$$ = yylex.(*PatitoLexer).Sem.CheckOperation($1, semantic.OpMult, $3)
	}
	| termino DIVIDE factor
	{
		$$ = yylex.(*PatitoLexer).Sem.CheckOperation($1, semantic.OpDiv, $3)
	}
	| factor
	{
		$$ = $1
	}
	;

factor:
	PAR_IZQ expresion PAR_DER
	{
		$$ = $2
	}
	| MAS factor %prec UPLUS
	{
		$$ = $2
	}
	| MENOS factor %prec UMINUS
	{
		$$ = $2
	}
	| ID
	{
		$$ = yylex.(*PatitoLexer).Sem.GetVarType($1)
	}
	| cte
	{
		$$ = $1
	}
	| llamada
	{
		$$ = semantic.TypeNula
	}
	;

cte:
	CTE_ENT
	{
		$$ = semantic.TypeEntero
	}
	| CTE_FLOT
	{
		$$ = semantic.TypeFlotante
	}
	;

%%

func Parse(l yyLexer) int {
	return yyParse(l)
}