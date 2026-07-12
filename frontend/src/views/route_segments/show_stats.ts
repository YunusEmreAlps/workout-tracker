import { html, LitElement, PropertyValues, TemplateResult } from "lit";
import { customElement, property } from "lit/decorators.js";
import {
  Chart,
  ScatterController,
  LineElement,
  PointElement,
  LinearScale,
  TimeScale,
  Tooltip,
  Legend,
} from "chart.js";
import "chartjs-adapter-date-fns";
import pluginTrendlineLinear from "chartjs-plugin-trendline";
import { localized } from "@lit/localize";
import { initLocalize } from "../../locale.js";

initLocalize();

interface TrendDataPoint {
  date: string;
  speed: string;
  type: string;
}

interface Translations {
  speedUnit: string;
  [key: string]: string;
}

@customElement("route-segment-stats")
@localized()
export class RouteSegmentStats extends LitElement {
  @property({
    converter: (v: string) => JSON.parse(v) as TrendDataPoint[],
  })
  data: TrendDataPoint[] = [];

  @property({
    attribute: "color-mode",
  })
  colorMode = "browser";

  @property()
  lang: string = null;

  @property({
    converter: (value: string) => JSON.parse(value) as Translations,
  })
  translations: Translations = null;

  private chart: Chart | null = null;

  public constructor() {
    super();
    Chart.register(
      ScatterController,
      LineElement,
      PointElement,
      LinearScale,
      TimeScale,
      Tooltip,
      Legend,
      pluginTrendlineLinear,
    );
  }

  private isDark(): boolean {
    if (this.colorMode === "dark") return true;
    if (this.colorMode === "light") return false;
    return window.matchMedia("(prefers-color-scheme: dark)").matches;
  }

  public override updated(_props: PropertyValues): void {
    super.updated(_props);

    if (!this.data || this.data.length === 0) return;

    const processedData = this.data.map((d) => ({
      ...d,
      type: d.type || "unknown",
    }));

    if (this.chart) {
      this.chart.destroy();
      this.chart = null;
    }

    const canvas = this.querySelector("canvas") as HTMLCanvasElement;
    if (!canvas) return;

    const speedUnit = this.translations?.speedUnit || "";
    const dark = this.isDark();
    const fgColor = dark ? "#e4e4e7" : "#27272a";
    const gridColor = dark ? "#3f3f46" : "#d4d4d8";

    const colorPalette = dark
      ? [
          { dot: "#60a5fa", trend: "#3b82f6" }, // Blue
          { dot: "#f87171", trend: "#ef4444" }, // Red
          { dot: "#34d399", trend: "#10b981" }, // Green
          { dot: "#fbbf24", trend: "#f59e0b" }, // Amber/Yellow
          { dot: "#c084fc", trend: "#a855f7" }, // Purple
          { dot: "#22d3ee", trend: "#06b6d4" }, // Cyan
          { dot: "#f472b6", trend: "#ec4899" }, // Pink
        ]
      : [
          { dot: "#1d4ed8", trend: "#2563eb" }, // Blue
          { dot: "#b91c1c", trend: "#dc2626" }, // Red
          { dot: "#047857", trend: "#10b981" }, // Green
          { dot: "#b45309", trend: "#d97706" }, // Amber/Yellow
          { dot: "#6b21a8", trend: "#8b5cf6" }, // Purple
          { dot: "#0e7490", trend: "#06b6d4" }, // Cyan
          { dot: "#be185d", trend: "#ec4899" }, // Pink
        ];

    const allTypes = Array.from(
      new Set(processedData.map((d) => d.type)),
    ).sort();
    const typeColorMap = new Map<string, { dot: string; trend: string }>();
    allTypes.forEach((type, index) => {
      typeColorMap.set(type, colorPalette[index % colorPalette.length]);
    });

    const datasets = allTypes.map((type) => {
      const typeData = processedData
        .filter((d) => d.type === type)
        .map((d) => ({
          x: new Date(d.date).valueOf(),
          y: parseFloat(d.speed),
        }))
        .filter((d) => !isNaN(d.y))
        .sort((a, b) => a.x - b.x);

      const defaultDotColor = dark ? "#fef3c7" : "#b45309";
      const defaultTrendColor = dark ? "#fbbf24" : "#d97706";
      const colors = typeColorMap.get(type) || {
        dot: defaultDotColor,
        trend: defaultTrendColor,
      };

      return {
        label: this.translations?.[type] || type,
        backgroundColor: colors.dot,
        borderColor: colors.dot,
        data: typeData,
        showLine: true,
        borderWidth: 2,
        pointRadius: 4,
        trendlineLinear: {
          colorMin: colors.trend,
          colorMax: colors.trend,
          width: 2,
          lineStyle: "solid",
        },
      };
    });

    this.chart = new Chart(canvas, {
      type: "scatter",
      data: {
        datasets: datasets as never[],
      },
      options: {
        maintainAspectRatio: false,
        animation: false,
        scales: {
          x: {
            type: "time",
            time: { unit: "month" },
            ticks: { color: fgColor },
            grid: { color: gridColor },
          },
          y: {
            ticks: {
              color: fgColor,
              callback: (val) => `${Number(val).toFixed(1)} ${speedUnit}`,
            },
            grid: { color: gridColor },
          },
        },
        plugins: {
          legend: {
            display: true,
            labels: { color: fgColor },
          },
          tooltip: {
            callbacks: {
              label: (item) => {
                const label = item.dataset.label || "";
                return `${label}: ${(item.raw as { y: number }).y.toFixed(
                  1,
                )} ${speedUnit}`;
              },
              title: (items) =>
                new Date(items[0].parsed.x).toLocaleDateString(),
            },
          },
        },
      },
    });
  }

  public render(): TemplateResult {
    this.style.display = "block";
    this.style.width = "100%";
    this.style.height = "100%";

    return html`
      <div class="h-[300px] md:h-[400px]">
        <canvas></canvas>
      </div>
    `;
  }

  protected createRenderRoot() {
    return this;
  }
}
