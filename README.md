# spot-cli

ferramenta em linha de comando para customizar o cliente do spotify.
suporte para windows, macos e linux.

**features:**

- alterar as cores de toda a ui
- injetar css para customização avançada
- injeção de extensões (script javascript) para estender funcionalidades, manipular ui e controlar o player
- habilitar algumas features adicionais, escondidas
- remover componentes para melhorar a performance

## instalação

1. baixe o pacote correto para seu sistema operacional: https://github.com/lifeondisplay/spot/releases
2. unpack

#### windows

no powershell, rode os seguintes comandos:

```powershell
Expand-Archive "$HOME\Downloads\spot-xxx.zip" "$HOME\spot"
```

com `$HOME/Downloads/spot-xxx.tar.gz` no path direto para o pacote baixado.

extraia o pacote zip

opcionalmente, rode:

```powershell
Add-Content $PROFILE "Set-Alias spot `"$HOME\spot\spot.exe`""
```

reinicie o powershell. agora você pode rodar `spot` em qualquer lugar.

#### linux e macos

no terminal, rode os seguintes comandos:

```bash
mkdir ~/spot
tar xzf ~/Downloads/spot-xxx.tar.gz -C ~/spot
```

com `~/Downloads/spot-xxx.tar.gz` apenas o path direto do pacote baixado

opcionalmente, rode:

```bash
sudo ln -s ~/spot/spot /usr/bin/spot
```

agora você pode rodar `spot` em qualquer lugar

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

para outros comandos e informação de flags adicionais, por favor rode:

```bash
spot --help
```

## customização

#### arquivo de configuração

localizado em:
**windows**: `%userprofile%\.spot\config.ini`
**linux:** `~/.spot/config.ini`  
**macos:** `~/spot_data/config.ini`  

para informação de detalhe de cada field de configuração, por favor rode:

```bash
spot --help config
```

#### temas

há 2 espaços onde você pode inserir seus temas:

1. a pasta `Themes` no diretório home
**windows**: `%userprofile%\.spot\Themes\`
**linux:** `~/.spot/Themes/`  
**macos:** `~/spot_data/Themes`    
2. a pasta `Themes` no diretório do executável spot

caso haja 2 temas contendo o mesmo nome, o tema no diretório home é o priorizado.

#### extensões

adicione os nomes da sua extensão desejada em config, separado pelo caractere `|`.

exemplo:

```ini
[AdditionalOptions]
...
extensions
autoSkipExplicit.js|queueAll.js|djMode.js|shuffle+.js|trashbin.js
```

arquivos de extensão podem ser armazenados em:

- na pasta `Extensions` no diretório home:
**windows:** `%userprofile%\.spot\Extensions\`
**linux:** `~/.spot/Extensions/`
**macos:** `~/spot_data/Extensions`
- na pasta `Extensions` no diretório do executável do spot

se houver 2 ramais com o mesmo nome, o ramal no diretório inicial será priorizado.

algumas apis do spotify vazaram e foram colocadas no objeto global `spot`. confira `global.d.ts` para documentação da api.

## desenvolvimento

### requisitos

- [go](https://go.dev/dl/)

clone o repositório e baixe as dependências:

```bash
go get github.com/lifeondisplay/spot
```

### build

#### windows

```powershell
cd $HOME\go\src\github.com\lifeondisplay\spot
go build -o spot.exe
```

#### linux e macos

```bash
cd ~/go/src/github.com/lifeondisplay/spot
go build -o spot
```

## em breve

- sass
- injetar aplicativos customizados
