FROM golang:1.13-alpine AS build
WORKDIR /social_media
COPY . .
RUN go build -o server .

FROM alpine:3.12
EXPOSE 9092
RUN mkdir /application
WORKDIR /application
COPY --from=build /social_media/frontend ./frontend/
COPY --from=build /social_media/robots.txt .
COPY --from=build /social_media/media ./media/
COPY --from=build /social_media/server .
RUN chmod +x ./server
CMD ["./server"]
