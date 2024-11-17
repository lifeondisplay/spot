# spot-cli

## instalação

1. baixe o pacote correto para seu sistema operacional: https://github.com/lifeondisplay/spot/releases
2. unpack

#### windows

extraia o pacote zip

para utilizar o spot, você pode rodar o `spot.exe` diretamente com o seu path,
ou adicione opcionalmente seu diretório no path de ambiente para rodar o `spot` aonde quiser

#### linux e macos

no terminal, rode os seguintes comandos:

```bash
cd ~/
mkdir spot
cd spot
tar xzf ~/Downloads/spot-xxx.tar.gz
```

com `~/Downloads/spot-xxx.tar.gz` apenas o path do pacote baixado

opcionalmente, rode:

```bash
echo 'spot=~/spot/spot' >> .bashrc
```

você pode rodar `spot` em qualquer lugar

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

## customização

#### arquivo de configuração

está localizado em:
**windows**: `%userprofile%\.spot\config.ini`
**linux e macos**: `~/.spot/config.ini`

#### temas

há 2 espaços onde você pode inserir seus temas:

1. a pasta `Themes` no diretório home
**windows**: `%userprofile%\.spot\Themes\`
**linux e macos**: `~/.spot/Themes/`
2. a pasta `Themes` no diretório do executável spot

caso haja 2 temas contendo o mesmo nome, o tema no diretório home é o priorizado.

## desenvolvimento

### requisitos

- [go](https://go.dev/dl/)

```bash
git clone https://github.com/lifeondisplay/spot
```

### build

```bash
cd spot-cli
go build src/spot.go
```

## em breve

- implementação de recursos futuros
- sass
- assistir a alteração dos arquivos de temas e aplicar automaticamente
- injetar extensões e aplicativos customizados
