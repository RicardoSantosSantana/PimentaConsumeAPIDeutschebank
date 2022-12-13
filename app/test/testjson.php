<?php

$data = `{"originIban":"DE00500700100200000885", "csv": { "name":"Ricardo" }}`;



$transaction = new stdClass;

$transaction->originIban = "";
$transaction->amount = "";
$transaction->counterPartyName = "";
$transaction->{'csv'} = json_decode($data);

$checkExistTransaction = json_encode($transaction);

echo $checkExistTransaction;
