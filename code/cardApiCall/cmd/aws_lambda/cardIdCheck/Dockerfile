FROM golang:1.24 AS build
ENV GOOS=linux
ENV GOARCH=arm64
WORKDIR /cardApiCall
# Copy dependencies list
COPY go.mod go.sum ./
RUN go mod download
# Copy cmd/aws_lambda/cardIdCheck
COPY cmd/aws_lambda/cardIdCheck/ ./cmd/aws_lambda/cardIdCheck/
RUN go build -tags lambda.norpc -o main ./cmd/aws_lambda/cardIdCheck/main.go
# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /cardApiCall/main ./main
ENTRYPOINT [ "./main" ]
