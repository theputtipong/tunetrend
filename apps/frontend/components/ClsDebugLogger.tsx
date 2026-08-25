"use client";

import { useEffect } from "react";
import { onCLS, type CLSMetric } from "web-vitals";

// TEMPORARY: หา DOM node ที่ทำให้เกิด Cumulative Layout Shift จริง (P75 = 0.46 บน production)
// ก่อนจะแก้อะไรแบบเดา ลบ component นี้ทิ้งหลังจากเจอ+แก้ตัวการแล้ว
function logShift(metric: CLSMetric) {
  for (const entry of metric.entries) {
    for (const source of entry.sources ?? []) {
      console.warn(`[CLS] value=${entry.value.toFixed(4)} node=`, source.node, {
        previousRect: source.previousRect,
        currentRect: source.currentRect,
      });

      if (source.node instanceof HTMLElement) {
        const el = source.node;
        const prevOutline = el.style.outline;
        el.style.outline = "3px solid red";
        setTimeout(() => {
          el.style.outline = prevOutline;
        }, 3000);
      }
    }
  }
}

export function ClsDebugLogger() {
  useEffect(() => {
    onCLS(logShift, { reportAllChanges: true });
  }, []);

  return null;
}
