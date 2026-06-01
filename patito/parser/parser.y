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
	PROGRAMA ID PUNTOCOMA 
	{
		yylex.(*PatitoLexer).Gen.GenerateMainGoto()
	}
	vars funcs INICIO 
	{
		yylex.(*PatitoLexer).Gen.FillMainGoto()
	}
	cuerpo FIN
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
		yylex.(*PatitoLexer).Sem.SetFunctionStart($2, yylex.(*PatitoLexer).Gen.NextQuad())
	}
	PAR_IZQ params PAR_DER LLAVE_IZQ vars cuerpo LLAVE_DER PUNTOCOMA
	{
		yylex.(*PatitoLexer).Gen.GenerateEndFunc()
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
		destAddr, varType := yylex.(*PatitoLexer).Sem.GetVarAddressAndType($1)

		if varType != semantic.TypeError && $3 != semantic.TypeError {
			yylex.(*PatitoLexer).Gen.GenerateAssignment(destAddr, varType)
		}
	}
	;

condicion:
	SI PAR_IZQ expresion PAR_DER
	{
		yylex.(*PatitoLexer).Gen.StartIf()
	}
	cuerpo sino_opc PUNTOCOMA
	;

sino_opc:
	SINO
	{
		yylex.(*PatitoLexer).Gen.ElseIf()
	}
	cuerpo
	{
		yylex.(*PatitoLexer).Gen.EndIfElse()
	}
	| /* empty */
	{
		yylex.(*PatitoLexer).Gen.EndIf()
	}
	;

ciclo:
	MIENTRAS
	{
		yylex.(*PatitoLexer).Gen.StartWhile()
	}
	PAR_IZQ expresion PAR_DER
	{
		yylex.(*PatitoLexer).Gen.WhileCondition()
	}
	HAZ cuerpo PUNTOCOMA
	{
		yylex.(*PatitoLexer).Gen.EndWhile()
	}
	;

llamada:
	ID 
	{
		yylex.(*PatitoLexer).Sem.StartCall($1)
		yylex.(*PatitoLexer).Gen.GenerateERA($1)
	}
	PAR_IZQ argumentos PAR_DER
	{
		startQuad := yylex.(*PatitoLexer).Sem.EndCall($1)
		yylex.(*PatitoLexer).Gen.GenerateGOSUB($1, startQuad)
	}
	;

argumentos:
    expresion
    {
        param, ok := yylex.(*PatitoLexer).Sem.GetCurrentParam()
        if ok {
            yylex.(*PatitoLexer).Gen.GenerateParam(param.Type, param.Address)
        }
    }
    argumentos_prima
    | /* empty */
    ;

argumentos_prima:
    COMA expresion
    {
        param, ok := yylex.(*PatitoLexer).Sem.GetCurrentParam()
        if ok {
            yylex.(*PatitoLexer).Gen.GenerateParam(param.Type, param.Address)
        }
    }
    argumentos_prima
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
	{
		value := yylex.(*PatitoLexer).Gen.PopOperandForPrint()
		yylex.(*PatitoLexer).Gen.GeneratePrint(value)
	}
	| LETRERO
	{		
		addr := yylex.(*PatitoLexer).Sem.AddConstant($1, semantic.TypeString)
		yylex.(*PatitoLexer).Gen.GeneratePrint(addr)
	}
	;

expresion:
	exp
	{
		$$ = $1
	}
	| exp MAYOR
	{
		yylex.(*PatitoLexer).Gen.PushOperator(">")
	}
	exp
	{
		yylex.(*PatitoLexer).Gen.GenerateBinaryOperation()
		$$ = yylex.(*PatitoLexer).Gen.PeekType()
	}
	| exp MENOR
	{
		yylex.(*PatitoLexer).Gen.PushOperator("<")
	}
	exp
	{
		yylex.(*PatitoLexer).Gen.GenerateBinaryOperation()
		$$ = yylex.(*PatitoLexer).Gen.PeekType()
	}
	| exp DIF
	{
		yylex.(*PatitoLexer).Gen.PushOperator("!=")
	}
	exp
	{
		yylex.(*PatitoLexer).Gen.GenerateBinaryOperation()
		$$ = yylex.(*PatitoLexer).Gen.PeekType()
	}
	| exp IGUAL
	{
		yylex.(*PatitoLexer).Gen.PushOperator("==")
	}
	exp
	{
		yylex.(*PatitoLexer).Gen.GenerateBinaryOperation()
		$$ = yylex.(*PatitoLexer).Gen.PeekType()
	}
	;

exp:
	exp MAS
	{
		yylex.(*PatitoLexer).Gen.PushOperator("+")
	}
	termino
	{
		yylex.(*PatitoLexer).Gen.GenerateBinaryOperation()
		$$ = yylex.(*PatitoLexer).Gen.PeekType()
	}
	| exp MENOS
	{
		yylex.(*PatitoLexer).Gen.PushOperator("-")
	}
	termino
	{
		yylex.(*PatitoLexer).Gen.GenerateBinaryOperation()
		$$ = yylex.(*PatitoLexer).Gen.PeekType()
	}
	| termino
	{
		$$ = $1
	}
	;


termino:
	termino MULT
	{
		yylex.(*PatitoLexer).Gen.PushOperator("*")
	}
	factor
	{
		yylex.(*PatitoLexer).Gen.GenerateBinaryOperation()
		$$ = yylex.(*PatitoLexer).Gen.PeekType()
	}
	| termino DIVIDE
	{
		yylex.(*PatitoLexer).Gen.PushOperator("/")
	}
	factor
	{
		yylex.(*PatitoLexer).Gen.GenerateBinaryOperation()
		$$ = yylex.(*PatitoLexer).Gen.PeekType()
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
		addr, t := yylex.(*PatitoLexer).Sem.GetVarAddressAndType($1)
		yylex.(*PatitoLexer).Gen.PushOperand(addr, t)
		$$ = t
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
		addr := yylex.(*PatitoLexer).Sem.AddConstant($1, semantic.TypeEntero)
		yylex.(*PatitoLexer).Gen.PushOperand(addr, semantic.TypeEntero)
		$$ = semantic.TypeEntero
	}
	| CTE_FLOT
	{		
		addr := yylex.(*PatitoLexer).Sem.AddConstant($1, semantic.TypeFlotante)
		yylex.(*PatitoLexer).Gen.PushOperand(addr, semantic.TypeFlotante)
		$$ = semantic.TypeFlotante
	}
	;

%%

func Parse(l yyLexer) int {
	return yyParse(l)
}