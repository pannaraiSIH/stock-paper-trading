"use client";

import { useEffect, useRef } from "react";
import {
  createChart,
  AreaSeries,
  ColorType,
  type IChartApi,
  type ISeriesApi,
  type UTCTimestamp,
} from "lightweight-charts";
import type { Candle } from "@/types";

const COLORS = {
  line: "#5b86f5",
  topArea: "rgba(91, 134, 245, 0.28)",
  bottomArea: "rgba(91, 134, 245, 0)",
  grid: "#33394a",
  text: "#9aa2b1",
};

function parseCandleTime(datetime: string): UTCTimestamp {
  const iso =
    datetime.length <= 10
      ? `${datetime}T00:00:00Z`
      : `${datetime.replace(" ", "T")}Z`;
  return Math.floor(new Date(iso).getTime() / 1000) as UTCTimestamp;
}

interface PriceChartProps {
  candles: Candle[];
  loading: boolean;
}

export function PriceChart({ candles, loading }: PriceChartProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<IChartApi | null>(null);
  const seriesRef = useRef<ISeriesApi<"Area"> | null>(null);

  useEffect(() => {
    if (!containerRef.current) return;

    const chart = createChart(containerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: "transparent" },
        textColor: COLORS.text,
        fontSize: 10,
      },
      grid: {
        vertLines: { visible: false },
        horzLines: { color: COLORS.grid },
      },
      rightPriceScale: { borderVisible: false },
      timeScale: {
        borderVisible: false,
        timeVisible: true,
        secondsVisible: false,
      },
      crosshair: {
        vertLine: { color: COLORS.grid, labelBackgroundColor: COLORS.grid },
        horzLine: { color: COLORS.grid, labelBackgroundColor: COLORS.grid },
      },
      handleScroll: false,
      handleScale: false,
    });

    const series = chart.addSeries(AreaSeries, {
      lineColor: COLORS.line,
      topColor: COLORS.topArea,
      bottomColor: COLORS.bottomArea,
      lineWidth: 2,
      priceLineVisible: false,
    });

    chartRef.current = chart;
    seriesRef.current = series;

    const resize = () => {
      if (!containerRef.current) return;
      chart.applyOptions({
        width: containerRef.current.clientWidth,
        height: containerRef.current.clientHeight,
      });
    };
    resize();

    const observer = new ResizeObserver(resize);
    observer.observe(containerRef.current);

    return () => {
      observer.disconnect();
      chart.remove();
      chartRef.current = null;
      seriesRef.current = null;
    };
  }, []);

  useEffect(() => {
    const series = seriesRef.current;
    if (!series) return;

    const data = candles
      .map((c) => ({
        time: parseCandleTime(c.datetime),
        value: Number(c.close),
      }))
      .filter((d) => Number.isFinite(d.value))
      .sort((a, b) => a.time - b.time);

    series.setData(data);
    chartRef.current?.timeScale().fitContent();
  }, [candles]);

  return (
    <div
      ref={containerRef}
      className={`chart chart-fade${loading ? " is-swapping" : ""}`}
    >
      {!loading && candles.length === 0 && (
        <div className="chart__empty">No price data available</div>
      )}
    </div>
  );
}
