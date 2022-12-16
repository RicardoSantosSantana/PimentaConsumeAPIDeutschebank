<?php




function amount(int $valor)
{

    return $valor > 0;
}

var_dump(amount(10));
var_dump(amount(-10));
