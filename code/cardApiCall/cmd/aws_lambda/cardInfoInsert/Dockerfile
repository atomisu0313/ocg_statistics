FROM golang:1.24 AS build
ENV GOOS=linux
ENV GOARCH=arm64
WORKDIR /cardApiCall
# Copy dependencies list
COPY go.mod go.sum ./
RUN go mod download
# Copy internal and cmd/aws_lambda/cardInfoInsert
COPY internal/ ./internal/
COPY cmd/aws_lambda/cardInfoInsert/ ./cmd/aws_lambda/cardInfoInsert/
RUN go build -tags lambda.norpc -o main ./cmd/aws_lambda/cardInfoInsert/main.go
# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /cardApiCall/main ./main
ENTRYPOINT [ "./main" ]