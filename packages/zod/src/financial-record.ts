import { z } from "zod";
import { ZDate, ZDateTime, ZDecimalString, ZHTTPError, ZUUID } from "./common.js";

export const ZRecordType = z.enum(["income", "expense"]);
export const ZTrendInterval = z.enum(["weekly", "monthly"]);

export const ZFinancialRecord = z.object({
  id: ZUUID,
  userId: ZUUID,
  amount: ZDecimalString,
  type: ZRecordType,
  category: z.string(),
  date: ZDateTime,
  notes: z.string().nullable().optional(),
  createdAt: ZDateTime,
  updatedAt: ZDateTime,
});

export const ZCreateFinancialRecordBody = z.object({
  amount: z.string(),
  type: ZRecordType,
  category: z.string().min(2).max(100),
  date: ZDate,
  notes: z.string().max(500).optional(),
});

export const ZUpdateFinancialRecordBody = z.object({
  amount: z.string().optional(),
  type: ZRecordType.optional(),
  category: z.string().min(2).max(100).optional(),
  date: ZDate.optional(),
  notes: z.string().max(500).optional(),
});

export const ZListFinancialRecordsQuery = z.object({
  type: ZRecordType.optional(),
  category: z.string().optional(),
  dateFrom: ZDate.optional(),
  dateTo: ZDate.optional(),
});

export const ZDashboardSummaryQuery = z.object({
  dateFrom: ZDate.optional(),
  dateTo: ZDate.optional(),
  trendInterval: ZTrendInterval.optional(),
  trendPeriods: z.coerce.number().int().min(1).max(24).optional(),
  recentLimit: z.coerce.number().int().min(1).max(20).optional(),
});

export const ZDashboardCategoryTotal = z.object({
  type: ZRecordType,
  category: z.string(),
  total: ZDecimalString,
});

export const ZDashboardTrendPoint = z.object({
  interval: ZTrendInterval,
  periodStart: ZDateTime,
  periodEnd: ZDateTime,
  income: ZDecimalString,
  expenses: ZDecimalString,
  netBalance: ZDecimalString,
});

export const ZDashboardSummary = z.object({
  dateFrom: ZDateTime.optional(),
  dateTo: ZDateTime.optional(),
  trendInterval: ZTrendInterval,
  totalIncome: ZDecimalString,
  totalExpenses: ZDecimalString,
  netBalance: ZDecimalString,
  categoryTotals: z.array(ZDashboardCategoryTotal),
  recentActivity: z.array(ZFinancialRecord),
  trends: z.array(ZDashboardTrendPoint),
});

export const ZFinancialRecordList = z.array(ZFinancialRecord);

export const ZFinancialErrorResponse = ZHTTPError;
