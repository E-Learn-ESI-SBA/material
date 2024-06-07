FROM golang:alpine3.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
RUN chmod +x main
# Create Sys user named materials , with group materials
RUN addgroup -S materials && adduser -S materials -G materials
RUN chown -R materials:materials /app
EXPOSE 8080
USER materials
CMD [ "./main" ]








# Set the Current Working Directory inside the container
