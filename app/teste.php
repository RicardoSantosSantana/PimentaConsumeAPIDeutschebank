<?php

class TestBankTransactions
{
    public function getAccounts()
    {
        $accounts = [
            'accounts' => [
                ['iban' => "iban001"],
                ['iban' => "iban002"]
            ]
        ];
        return json_decode(json_encode($accounts));
    }
    public function call($url, $params)
    {
        $iban = $params['iban'];
        $offset = $params['offset'];
        $limit = $params['limit'];
        $send = $url . "?iban=$iban&offset=$offset&limit&$limit";
        echo $send . PHP_EOL;

        $retorno = new stdClass;
        $retorno->transactions = [];
        $retorno->totalItems = 500;
        return $retorno;
    }
    function saveTransactions($transactions, $iban)
    {
        echo "Salvando as transações." . PHP_EOL;
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

        $offset = $this->OffSet($response->totalItems, $limit);

        for ($i = 0; $i < $offset; $i++) {

            $response = $this->call($url, [
                'iban' => $iban,
                'offset' => $i,
                'limit' => $limit
            ]);

            $this->saveTransactions($response->transactions, $iban);
            $this->logInfo("Total items: " . $response->totalItems . " | limit: $limit | offset: $i");
        }
        echo "*********************************\n";
    }

    public function getTransactions($offset = 0, $accounts = null)
    {
        if (!$accounts) {
            $accounts = $this->getAccounts();
        }

        $this->logInfo("Starting importing from Deutschebank");
        echo "*********************************\n";

        foreach ($accounts->accounts as $account) {

            $this->getCallTransactions($account->iban);
        }
        $this->logInfo("End importing from Deutschebank");
    }
}

$transaction = new TestBankTransactions;
$transaction->getTransactions();
