import { initContract } from "@ts-rest/core";
import {
  ZCreateFinancialRecordBody,
  ZDashboardSummary,
  ZDashboardSummaryQuery,
  ZFinancialRecord,
  ZFinancialRecordList,
  ZHTTPError,
  ZListFinancialRecordsQuery,
  ZUpdateFinancialRecordBody,
  ZUUID,
} from "@boilerplate/zod";
import { z } from "zod";
import { getSecurityMetadata } from "@/utils.js";

const c = initContract();

export const financialRecordsContract = c.router({
  listFinancialRecords: {
    summary: "List financial records",
    path: "/api/v1/records",
    method: "GET",
    description: "Lists financial records for the authenticated user.",
    metadata: getSecurityMetadata(),
    query: ZListFinancialRecordsQuery,
    responses: {
      200: ZFinancialRecordList,
      400: ZHTTPError,
      401: ZHTTPError,
      403: ZHTTPError,
    },
  },
  getFinancialRecordById: {
    summary: "Get financial record",
    path: "/api/v1/records/:id",
    method: "GET",
    description: "Returns a financial record by UUID.",
    metadata: getSecurityMetadata(),
    pathParams: z.object({
      id: ZUUID,
    }),
    responses: {
      200: ZFinancialRecord,
      401: ZHTTPError,
      403: ZHTTPError,
      404: ZHTTPError,
    },
  },
  createFinancialRecord: {
    summary: "Create financial record",
    path: "/api/v1/records",
    method: "POST",
    description: "Creates a new financial record. Analyst and admin only.",
    metadata: getSecurityMetadata(),
    body: ZCreateFinancialRecordBody,
    responses: {
      201: ZFinancialRecord,
      400: ZHTTPError,
      401: ZHTTPError,
      403: ZHTTPError,
    },
  },
  updateFinancialRecord: {
    summary: "Update financial record",
    path: "/api/v1/records/:id",
    method: "PATCH",
    description: "Updates a financial record. Analyst and admin only.",
    metadata: getSecurityMetadata(),
    pathParams: z.object({
      id: ZUUID,
    }),
    body: ZUpdateFinancialRecordBody,
    responses: {
      200: ZFinancialRecord,
      400: ZHTTPError,
      401: ZHTTPError,
      403: ZHTTPError,
      404: ZHTTPError,
    },
  },
  deleteFinancialRecord: {
    summary: "Delete financial record",
    path: "/api/v1/records/:id",
    method: "DELETE",
    description: "Deletes a financial record. Admin only.",
    metadata: getSecurityMetadata(),
    pathParams: z.object({
      id: ZUUID,
    }),
    responses: {
      204: z.null(),
      401: ZHTTPError,
      403: ZHTTPError,
      404: ZHTTPError,
    },
  },
  getDashboardSummary: {
    summary: "Get dashboard summary",
    path: "/api/v1/dashboard/summary",
    method: "GET",
    description:
      "Returns summary totals, recent activity, category totals, and trend data for the authenticated user.",
    metadata: getSecurityMetadata(),
    query: ZDashboardSummaryQuery,
    responses: {
      200: ZDashboardSummary,
      400: ZHTTPError,
      401: ZHTTPError,
      403: ZHTTPError,
    },
  },
});
