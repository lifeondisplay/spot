# spot-cli

## requisitos

- [go](https://golang.org/dl/)

## clone

```bash
git clone https://github.com/lifeondisplay/spot-cli
```

## build

```bash
cd spot-cli

go build src/spot.go
```

## uso

rode uma vez sem comando para gerar o arquivo de configuração

```bash
spot
```

e então:

```bash
spot backup
```

e finalmente:

```bash
spot apply
```

depois de alterar o tema de cor e o css, rode `apply` novamente

## em breve

- sass
- assistir a alteração dos arquivos de temas e aplicar automaticamente
- injetar extensões e aplicativos customizados
