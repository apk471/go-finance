import { initContract } from "@ts-rest/core";
import { ZHTTPError, ZHealthResponse } from "@boilerplate/zod";

const c = initContract();

export const systemContract = c.router({
  getHealth: {
    summary: "Get service health",
    path: "/status",
    method: "GET",
    description: "Returns service health plus dependency checks.",
    responses: {
      200: ZHealthResponse,
      503: ZHealthResponse,
    },
  },
  getDocs: {
    summary: "Open documentation UI",
    path: "/docs",
    method: "GET",
    description: "Serves the HTML API documentation page.",
    responses: {
      200: c.otherResponse({
        contentType: "text/html",
        body: c.type<string>(),
      }),
      500: ZHTTPError,
    },
  },
  getOpenApiSpec: {
    summary: "Get OpenAPI JSON",
    path: "/static/openapi.json",
    method: "GET",
    description: "Returns the generated OpenAPI document consumed by the docs UI.",
    responses: {
      200: c.otherResponse({
        contentType: "application/json",
        body: c.type<Record<string, unknown>>(),
      }),
    },
  },
});
