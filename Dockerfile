FROM golang:1.27-alpine

# es recomendable instalar git, ya que algunas dependencias de Go pueden requerirlo para descargarse correctamente.
RUN go install github.com/air-verse/air@v1.61.7

#  el directorio de trabajo dentro del contenedor se establece en /app, lo que significa que cualquier comando que se ejecute después de esta línea se ejecutará en ese directorio.
WORKDIR /app

#esto es clave para el cacheo de dependencias, ya que si no se copia el go.mod y go.sum antes de copiar el resto del código, cada vez que se haga un cambio en el código, se volverán a descargar todas las dependencias, lo cual es innecesario y lento.
COPY go.mod go.sum ./

#El comando "go mod download" descarga todas las dependencias necesarias para la aplicación, según lo especificado en el archivo go.mod.
RUN go mod download

# copiamos el resto del código de la aplicación de source a destino (/app)
COPY . .

# el comando CMD indica el comando que se ejecutará cuando se inicie el contenedor. En este caso, se está indicando que se ejecute el comando "air", que es un servidor de desarrollo para aplicaciones Go.
CMD [ "air" ]
