<?php

class TestBankTransactions
{
    public function getAccounts()
    {
        $accounts = [
            'accounts' => [
                ['iban' => "DE00500700100200000867"],
                ['iban' => "DE00500700100200000123"]
            ]
        ];
        return json_decode(json_encode($accounts));
    }

    public function call($url, $params)
    {
        $iban = $params['iban'];
        $offset = $params['offset'];
        $limit = $params['limit'];

        $retorno = new stdClass;
        $retorno->transactions = [];

        if ($iban == "DE00500700100200000867") {
            $retorno->totalItems = 500;
        } else {
            $retorno->totalItems = 1250;
        }

        return $retorno;
    }

    function saveTransactions($transactions, $iban)
    {
        echo "  Salvando as transações." . PHP_EOL;
    }

    public function logInfo($msg)
    {
        echo $msg . PHP_EOL;
    }

    public function OffSet($totalItems, $limit)
    {
        $offset =  $totalItems / $limit;

        if (is_float($offset)) {
            return (int)$offset + 1;
        }

        return (int)$offset;
    }

    public function getCallTransactions($iban = null, $limit = 200, $url = 'banking/transactions/v2')
    {
        // chama a primeira vez apenas para pegar o total de itens.
        $response = $this->call($url, ['iban' => $iban, 'offset' => 0, 'limit' => $limit]);

        $TotalOffset = $this->OffSet($response->totalItems, $limit);
        echo "\n -|- IBAN = $iban \n";
        for ($offset = 0; $offset < $TotalOffset; $offset++) {

            $response = $this->call($url, [
                'iban' => $iban,
                'offset' => $offset,
                'limit' => $limit
            ]);
            echo "\n  OFFSET = $offset \n";
            $this->saveTransactions($response->transactions, $iban);
            $this->logInfo("  Total items: " . $response->totalItems . " | limit: $limit | offset: $offset");
            $send = $url . "?iban=$iban&offset=$offset&limit&$limit";
            echo "  Api: " . $send . PHP_EOL;
        }

        echo "___________________________________________________________________________________\n";
    }

    public function getTransactions($offset = 0, $accounts = null)
    {
        if (!$accounts) {
            $accounts = $this->getAccounts();
        }

        $this->logInfo("\nStarting importing from Deutschebank\n");
        echo "---------------------------------------------------------\n";

        foreach ($accounts->accounts as $account) {

            $this->getCallTransactions($account->iban);
        }

        $this->logInfo("\nEnd importing from Deutschebank\n");
    }
}

$transaction = new TestBankTransactions;
$transaction->getTransactions();
