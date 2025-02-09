FROM gcr.io/distroless/cc AS base
COPY apps/eval/.out/twir-eval /app/twir-eval
ENV NODE_ENV=production
CMD ["/app/twir-eval"]
