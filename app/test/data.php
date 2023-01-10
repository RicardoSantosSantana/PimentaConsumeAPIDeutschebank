<?php

//&bookingDateFrom=2022-05-14&bookingDateTo=2022-05-14

$bookingDateFrom = date("Y-m-d", strtotime("-31 days"));
$bookingDateTo = date("Y-m-d");
 

//OutPut
echo $bookingDateTo;
echo PHP_EOL;
echo $bookingDateFrom;