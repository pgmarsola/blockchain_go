# Blockchain Go

#### Projeto criado para aprofundar os estudos em blockchain.

Este é um projeto de blockchain local desenvolvimento em Go baseado em Ethereum, tem como objetivo auxiliar nos estudos relacionados a blockchain.

## 💻 Pré-requisitos

Antes de começar, certifique-se que possui todos os requisitos necessários para executar o projeto

- `< Go / v1.23.2 >`

## 🚀 Executando

Para executar o **blockchain_go**, siga estas etapas:

Se na execução falhar por não conseguir baixar as dependencias automaticamente, execute o comando abaixo:

```
  go get < nome dependencia >
```

- Iniciar o blockchain

```
  go run main.go createblockchain -address < string >
```

> [!NOTE]
> Ao iniciar o blockchain, inicia-se junto o minerador, onde ficará printando no seu terminal todos os blocos minerados, isso torna seu terminal inutilizavel, para executar os demais comandos, deve-se abrir um novo terminal.

- Criar carteira

```
  go run main.go createwallet
```

- Listar todos os endereços de carteiras criados de forma criptografada

```
  go run main.go listadresses
```

- Consultar saldo em carteira

```
  go run main.go getbalance -address < wallet address >
```

- Realizar uma transação

```
  go run main.go send -from < wallet address > -to < wallet address > -amount < value >
```

- Consultar listagem de blocos criados

```
  go run main.go printchain
```

- Lista de comandos

```
  go run main.go
```

## 🤝 Contribuindo

Este projeto não aceita contribuições no momento.
