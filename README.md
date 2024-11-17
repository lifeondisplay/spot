# spot-cli

ferramenta em linha de comando para customizar o cliente do spotify.
suporte para windows, macos e linux.

**features:**

- alterar as cores de toda a ui
- injetar css para customização avançada
- habilitar algumas features adicionais, escondidas
- remover componentes para melhorar a performance

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

## uso básico

rode uma vez sem comando para gerar o arquivo de configuração

```bash
spot
```

certifique-se de que o arquivo de configuração foi criado com sucesso e não há erro, então execute:

```bash
spot backup apply enable-devtool
```

e a partir de agora, após alterar as cores em `color.ini` ou no css em `user.css`, você só precisa executar:

```bash
spot update
```

para atualizar o seu tema.

no spotify, pressione <kbd>ctrl</kbd> + <kbd>shift</kbd> + <kbd>r</kbd>/<kbd>command</kbd> + <kbd>shift</kbd> + <kbd>r</kbd> para recarregar e receber atualização visual do seu tema.

## customização

#### arquivo de configuração

está localizado em:
**windows**: `%userprofile%\.spot\config.ini`
**linux:** `~/.spot/config.ini`  
**macos:** `~/spot_data/config.ini`  

#### temas

há 2 espaços onde você pode inserir seus temas:

1. a pasta `Themes` no diretório home
**windows**: `%userprofile%\.spot\Themes\`
**linux** `~/.spot/Themes/`  
**macos:** `~/spot_data/Themes`    
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

- sass
- injetar extensões e aplicativos customizados
