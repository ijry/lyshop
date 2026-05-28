<template>
  <view class="ly-charts-candlestick" :style="{ width: containerWidth, height: containerHeight }">
    <ly-canvas
      ref="canvasRef"
      class="chart-canvas"
      :canvas-id="canvasId"
      width="100"
      height="100"
      :use-root-height-and-width="true"
      @ready="handleCanvasReady"
      @touchstart="handleTouchStart"
      @touchmove="handleTouchMove"
      @touchend="handleTouchEnd"
    />
  </view>
</template>

<script>
import chartHelper from '../../libs/util/chartHelper.js';

export default {
  name: 'ly-charts-candlestick',
  props: {
    option: {
      type: Object,
      default: () => ({})
    },
    height: {
      type: [String, Number],
      default: 420
    },
    width: {
      type: [String, Number],
      default: '100%'
    }
  },
  emits: ['click', 'zoom', 'tooltipShow'],
  data() {
    return {
      canvasId: `candlestick-chart-${Date.now()}`,
      ctx: null,
      canvasWidth: 0,
      canvasHeight: 0,
      chartGrid: { top: 20, right: 20, bottom: 50, left: 56, width: 0, height: 0 },
      priceGrid: null,
      volumeGrid: null,
      candleRegions: [],
      volumeRegions: [],
      markPointRegions: [],
      sliderRect: null,
      zoomStart: 0,
      zoomEnd: 100,
      renderState: null,
      activePointer: null,
      touchInfo: {
        startX: 0,
        startY: 0,
        lastX: 0,
        lastY: 0,
        draggingSlider: false,
        sliderDragMode: '',
        originZoomStart: 0,
        originZoomEnd: 100
      }
    };
  },
  computed: {
    containerHeight() {
      return typeof this.height === 'number' ? `${this.height}px` : this.height;
    },
    containerWidth() {
      return typeof this.width === 'number' ? `${this.width}px` : this.width;
    }
  },
  watch: {
    option: {
      handler() {
        this.syncZoomFromOption();
        this.drawChart();
      },
      deep: true
    }
  },
  methods: {
    initCanvas() {
      const canvasRef = this.$refs.canvasRef;
      if (canvasRef && typeof canvasRef.refresh === 'function') {
        canvasRef.refresh();
      }
    },
    handleCanvasReady(event) {
      const canvasRef = this.$refs.canvasRef;
      this.ctx = canvasRef;
      this.canvasWidth = event.width || canvasRef.getWidth();
      this.canvasHeight = event.height || canvasRef.getHeight();
      this.syncZoomFromOption();
      this.drawChart();
    },
    syncZoomFromOption() {
      const zoom = Array.isArray(this.option?.dataZoom) && this.option.dataZoom.length > 0
        ? this.option.dataZoom[0]
        : {};
      const start = zoom.start !== undefined ? Number(zoom.start) : 0;
      const end = zoom.end !== undefined ? Number(zoom.end) : 100;
      this.zoomStart = Math.max(0, Math.min(100, start));
      this.zoomEnd = Math.max(this.zoomStart + 1, Math.min(100, end));
    },
    parseLayoutValue(value, total, fallback) {
      if (value === undefined || value === null || value === 'auto') {
        return fallback;
      }
      if (typeof value === 'number') {
        return value;
      }
      if (typeof value !== 'string') {
        return fallback;
      }
      if (value.indexOf('%') !== -1) {
        return total * parseFloat(value) / 100;
      }
      if (/(rpx|upx)$/.test(value)) {
        return uni.upx2px(parseInt(value, 10));
      }
      return parseFloat(value) || fallback;
    },
    parseNumber(value) {
      if (value === '-' || value === null || value === undefined || value === '') {
        return null;
      }
      const numberValue = Number(value);
      return Number.isNaN(numberValue) ? null : numberValue;
    },
    normalizeSeriesValue(rawValue) {
      if (rawValue && typeof rawValue === 'object' && !Array.isArray(rawValue)) {
        return {
          value: this.parseNumber(rawValue.value),
          itemStyle: rawValue.itemStyle || {},
          raw: rawValue
        };
      }
      return {
        value: this.parseNumber(rawValue),
        itemStyle: {},
        raw: rawValue
      };
    },
    applyOpacity(color, opacity = 1) {
      if (opacity >= 1 || typeof color !== 'string') {
        return color;
      }
      if (color.indexOf('#') === 0) {
        const hex = color.replace('#', '');
        const safeHex = hex.length === 3
          ? `${hex[0]}${hex[0]}${hex[1]}${hex[1]}${hex[2]}${hex[2]}`
          : hex;
        const r = parseInt(safeHex.substring(0, 2), 16);
        const g = parseInt(safeHex.substring(2, 4), 16);
        const b = parseInt(safeHex.substring(4, 6), 16);
        return `rgba(${r}, ${g}, ${b}, ${opacity})`;
      }
      if (color.indexOf('rgb(') === 0) {
        return color.replace('rgb(', 'rgba(').replace(')', `, ${opacity})`);
      }
      return color;
    },
    formatNumber(value, digits = 2) {
      const numberValue = this.parseNumber(value);
      if (numberValue === null) {
        return '--';
      }
      return Number(numberValue).toFixed(digits);
    },
    formatCompactValue(value, digits = 2) {
      const numberValue = this.parseNumber(value);
      if (numberValue === null) {
        return '--';
      }
      const absValue = Math.abs(numberValue);
      if (absValue >= 100000000) {
        return `${(numberValue / 100000000).toFixed(digits)}亿`;
      }
      if (absValue >= 10000) {
        return `${(numberValue / 10000).toFixed(digits)}万`;
      }
      if (absValue >= 1000) {
        return `${(numberValue / 1000).toFixed(digits)}K`;
      }
      return Number(numberValue).toFixed(digits);
    },
    formatAxisLabel(value, axisLabel = {}, digits = 2, compact = false) {
      if (typeof axisLabel.formatter === 'function') {
        try {
          return String(axisLabel.formatter(value));
        } catch (error) {
          return compact ? this.formatCompactValue(value, digits) : this.formatNumber(value, digits);
        }
      }
      if (typeof axisLabel.formatter === 'string') {
        return axisLabel.formatter.replace('{value}', compact ? this.formatCompactValue(value, digits) : this.formatNumber(value, digits));
      }
      return compact ? this.formatCompactValue(value, digits) : this.formatNumber(value, digits);
    },
    truncateText(text, maxLength = 14) {
      const safeText = String(text ?? '');
      if (safeText.length <= maxLength) {
        return safeText;
      }
      return `${safeText.slice(0, maxLength - 1)}…`;
    },
    measureTextWidth(text, fontSize = 12) {
      if (!this.ctx || typeof this.ctx.measureText !== 'function') {
        return String(text).length * fontSize * 0.6;
      }
      this.ctx.setFontSize(fontSize);
      return this.ctx.measureText(String(text)).width || (String(text).length * fontSize * 0.6);
    },
    getAxisOption(source, index = 0) {
      if (Array.isArray(source)) {
        return source[index] || source[0] || {};
      }
      return source || {};
    },
    getSeriesCollection() {
      const series = Array.isArray(this.option?.series) ? this.option.series : [];
      let candlestickSeries = null;
      const lineSeries = [];
      const barSeries = [];
      series.forEach((serie) => {
        if (serie.type === 'candlestick' && !candlestickSeries) {
          candlestickSeries = serie;
        } else if (serie.type === 'line') {
          lineSeries.push(serie);
        } else if (serie.type === 'bar') {
          barSeries.push(serie);
        }
      });
      return { candlestickSeries, lineSeries, barSeries, allSeries: series };
    },
    getLegendColors(allSeries, legendData) {
      return legendData.map((name, index) => {
        const matched = allSeries.find((item) => item.name === name) || allSeries[index] || {};
        if (matched.type === 'candlestick') {
          return matched.itemStyle?.color || '#ec0000';
        }
        return matched.color || matched.lineStyle?.color || matched.itemStyle?.color || chartHelper.getColor(index + 1);
      });
    },
    drawTitle() {
      const title = this.option?.title;
      if (!title || !title.text) {
        return 0;
      }

      const fontSize = title.textStyle?.fontSize || 16;
      const color = title.textStyle?.color || '#333';
      const subText = title.subtext || '';
      const subFontSize = title.subtextStyle?.fontSize || 12;
      const subColor = title.subtextStyle?.color || '#666';
      let align = 'center';
      let x = this.canvasWidth / 2;
      if (title.left === 'left' || title.left === 0 || title.left === '0') {
        align = 'left';
        x = 10;
      } else if (title.left === 'right') {
        align = 'right';
        x = this.canvasWidth - 10;
      } else if (typeof title.left === 'number' || (typeof title.left === 'string' && title.left.indexOf('%') !== -1)) {
        align = 'left';
        x = this.parseLayoutValue(title.left, this.canvasWidth, 10);
      }

      this.ctx.setFontSize(fontSize);
      this.ctx.setFillStyle(color);
      this.ctx.setTextAlign(align);
      this.ctx.setTextBaseline('top');
      this.ctx.fillText(title.text, x, 10);

      let height = fontSize + 12;
      if (subText) {
        this.ctx.setFontSize(subFontSize);
        this.ctx.setFillStyle(subColor);
        this.ctx.fillText(subText, x, 10 + fontSize + 4);
        height += subFontSize + 6;
      }
      return height;
    },
    createBaseGrid(titleHeight, showSlider) {
      const gridOptions = Array.isArray(this.option?.grid)
        ? this.option.grid
        : [this.option?.grid || {}];
      const firstGrid = gridOptions[0] || {};
      const grid = {
        left: this.parseLayoutValue(firstGrid.left, this.canvasWidth, 56),
        right: this.parseLayoutValue(firstGrid.right, this.canvasWidth, 20),
        top: this.parseLayoutValue(firstGrid.top, this.canvasHeight, Math.max(20, titleHeight + 10)),
        bottom: this.parseLayoutValue(firstGrid.bottom, this.canvasHeight, showSlider ? 64 : 38)
      };
      grid.width = Math.max(0, this.canvasWidth - grid.left - grid.right);
      grid.height = Math.max(0, this.canvasHeight - grid.top - grid.bottom);
      return grid;
    },
    normalizeGridOption(gridOption, fallbackGrid) {
      const left = this.parseLayoutValue(gridOption?.left, this.canvasWidth, fallbackGrid.left);
      const right = this.parseLayoutValue(gridOption?.right, this.canvasWidth, fallbackGrid.right);
      const top = this.parseLayoutValue(gridOption?.top, this.canvasHeight, fallbackGrid.top);
      let bottom = this.parseLayoutValue(gridOption?.bottom, this.canvasHeight, fallbackGrid.bottom);
      let height = gridOption?.height !== undefined
        ? this.parseLayoutValue(gridOption.height, this.canvasHeight, fallbackGrid.height)
        : Math.max(0, this.canvasHeight - top - bottom);

      if (gridOption?.height !== undefined && gridOption?.bottom === undefined) {
        bottom = Math.max(0, this.canvasHeight - top - height);
      }
      if (gridOption?.bottom !== undefined && gridOption?.height === undefined) {
        height = Math.max(0, this.canvasHeight - top - bottom);
      }

      return {
        left,
        right,
        top,
        bottom,
        width: Math.max(0, this.canvasWidth - left - right),
        height: Math.max(0, height)
      };
    },
    createPanelLayout(baseGrid, hasVolume) {
      const gridOptions = Array.isArray(this.option?.grid)
        ? this.option.grid
        : [this.option?.grid || {}];

      if (gridOptions.length > 1) {
        const priceGrid = this.normalizeGridOption(gridOptions[0], baseGrid);
        const fallbackVolume = {
          left: priceGrid.left,
          right: priceGrid.right,
          top: priceGrid.top + priceGrid.height + 12,
          bottom: baseGrid.bottom,
          width: priceGrid.width,
          height: Math.max(48, Math.min(100, baseGrid.height * 0.24))
        };
        const volumeGrid = hasVolume ? this.normalizeGridOption(gridOptions[1] || {}, fallbackVolume) : null;
        return { priceGrid, volumeGrid };
      }

      if (!hasVolume) {
        return {
          priceGrid: {
            ...baseGrid
          },
          volumeGrid: null
        };
      }

      const panelGap = this.parseLayoutValue(this.option?.panelGap, this.canvasHeight, 12);
      const maxVolumeHeight = Math.max(56, Math.min(120, baseGrid.height * 0.34));
      const defaultVolumeHeight = Math.max(56, Math.min(maxVolumeHeight, baseGrid.height * 0.24));
      const volumeHeight = Math.max(
        56,
        Math.min(
          maxVolumeHeight,
          this.parseLayoutValue(this.option?.volumeHeight, baseGrid.height, defaultVolumeHeight)
        )
      );
      const priceHeight = Math.max(80, baseGrid.height - volumeHeight - panelGap);

      const priceGrid = {
        left: baseGrid.left,
        right: baseGrid.right,
        top: baseGrid.top,
        bottom: this.canvasHeight - baseGrid.top - priceHeight,
        width: baseGrid.width,
        height: priceHeight
      };
      const volumeGrid = {
        left: baseGrid.left,
        right: baseGrid.right,
        top: priceGrid.top + priceGrid.height + panelGap,
        bottom: baseGrid.bottom,
        width: baseGrid.width,
        height: Math.max(40, this.canvasHeight - (priceGrid.top + priceGrid.height + panelGap) - baseGrid.bottom)
      };
      return { priceGrid, volumeGrid };
    },
    getVisibleRange(totalCount) {
      if (totalCount <= 0) {
        return { start: 0, end: -1, count: 0 };
      }
      const start = Math.max(0, Math.floor(totalCount * this.zoomStart / 100));
      const end = Math.min(totalCount - 1, Math.ceil(totalCount * this.zoomEnd / 100) - 1);
      return {
        start,
        end: Math.max(start, end),
        count: Math.max(1, Math.max(start, end) - start + 1)
      };
    },
    buildVisibleCandles(candlestickSeries, categories) {
      const rawData = Array.isArray(candlestickSeries?.data) ? candlestickSeries.data : [];
      const totalCount = Math.min(categories.length || rawData.length, rawData.length);
      const visibleRange = this.getVisibleRange(totalCount);
      const items = [];
      for (let index = visibleRange.start; index <= visibleRange.end; index++) {
        const rawItem = rawData[index];
        if (!Array.isArray(rawItem) || rawItem.length < 4) {
          continue;
        }
        items.push({
          index,
          localIndex: items.length,
          category: categories[index] !== undefined ? categories[index] : `${index}`,
          open: Number(rawItem[0]),
          close: Number(rawItem[1]),
          low: Number(rawItem[2]),
          high: Number(rawItem[3]),
          raw: rawItem
        });
      }
      return {
        items,
        totalCount,
        visibleRange
      };
    },
    calculatePriceRange(visibleCandles, lineSeries, yAxisOption) {
      let min = Number.MAX_VALUE;
      let max = Number.MIN_VALUE;
      let hasData = false;

      visibleCandles.items.forEach((item) => {
        min = Math.min(min, item.low, item.open, item.close);
        max = Math.max(max, item.high, item.open, item.close);
        hasData = true;
      });

      lineSeries.forEach((serie) => {
        const data = Array.isArray(serie.data) ? serie.data : [];
        visibleCandles.items.forEach((item) => {
          const value = this.parseNumber(data[item.index]);
          if (value !== null) {
            min = Math.min(min, value);
            max = Math.max(max, value);
            hasData = true;
          }
        });
      });

      if (!hasData) {
        return { min: 0, max: 1 };
      }

      if (typeof yAxisOption?.min === 'number') {
        min = yAxisOption.min;
      }
      if (typeof yAxisOption?.max === 'number') {
        max = yAxisOption.max;
      }

      if (max === min) {
        return { min: min - 1, max: max + 1 };
      }
      const padding = (max - min) * 0.05;
      return {
        min: typeof yAxisOption?.min === 'number' ? min : min - padding,
        max: typeof yAxisOption?.max === 'number' ? max : max + padding
      };
    },
    calculateVolumeRange(visibleCandles, barSeries, yAxisOption) {
      let min = Number.MAX_VALUE;
      let max = Number.MIN_VALUE;
      let hasData = false;

      barSeries.forEach((serie) => {
        const data = Array.isArray(serie.data) ? serie.data : [];
        visibleCandles.items.forEach((item) => {
          const value = this.normalizeSeriesValue(data[item.index]).value;
          if (value !== null) {
            min = Math.min(min, value);
            max = Math.max(max, value);
            hasData = true;
          }
        });
      });

      if (!hasData) {
        return { min: 0, max: 1, baseline: 0 };
      }

      if (typeof yAxisOption?.min === 'number') {
        min = yAxisOption.min;
      } else if (min > 0) {
        min = 0;
      }

      if (typeof yAxisOption?.max === 'number') {
        max = yAxisOption.max;
      } else {
        max += (max - min || 1) * 0.08;
      }

      if (max === min) {
        max = min + 1;
      }

      return {
        min,
        max,
        baseline: Math.min(Math.max(0, min), max)
      };
    },
    getCategoryCenter(localIndex, count, grid) {
      return grid.left + ((localIndex + 0.5) / Math.max(count, 1)) * grid.width;
    },
    valueToY(value, min, max, grid) {
      return grid.top + grid.height - ((value - min) / Math.max(max - min, 1e-9)) * grid.height;
    },
    yToValue(y, min, max, grid) {
      const ratio = (grid.top + grid.height - y) / Math.max(grid.height, 1);
      return min + Math.max(0, Math.min(1, ratio)) * (max - min);
    },
    getValueByDim(item, dim = 'close') {
      if (dim === 'open') return item.open;
      if (dim === 'highest' || dim === 'high') return item.high;
      if (dim === 'lowest' || dim === 'low') return item.low;
      return item.close;
    },
    getSeriesColor(serie, index, fallback) {
      return serie.color || serie.lineStyle?.color || serie.itemStyle?.color || fallback || chartHelper.getColor(index + 1);
    },
    drawYAxisSplitArea(ticks, grid, yAxisOption) {
      const splitArea = yAxisOption?.splitArea;
      if (!splitArea?.show) {
        return;
      }
      const colors = Array.isArray(splitArea.areaStyle?.color) && splitArea.areaStyle.color.length > 0
        ? splitArea.areaStyle.color
        : ['#f8fafc', '#ffffff'];

      for (let i = 0; i < ticks.tickCount; i++) {
        const topValue = ticks.min + (i + 1) * ticks.step;
        const bottomValue = ticks.min + i * ticks.step;
        const topY = this.valueToY(topValue, ticks.min, ticks.max, grid);
        const bottomY = this.valueToY(bottomValue, ticks.min, ticks.max, grid);
        this.ctx.setFillStyle(colors[i % colors.length]);
        this.ctx.fillRect(grid.left, topY, grid.width, bottomY - topY);
      }
    },
    drawYAxis(ticks, grid, yAxisOption, compact = false) {
      const axisLabel = yAxisOption?.axisLabel || {};
      const splitLine = yAxisOption?.splitLine || {};
      const lineColor = splitLine.lineStyle?.color || '#e5e7eb';
      const labelColor = axisLabel.color || '#666';
      const digits = compact ? 1 : 2;
      this.ctx.setFontSize(axisLabel.fontSize || 11);
      this.ctx.setTextAlign('right');
      this.ctx.setTextBaseline('middle');

      for (let i = 0; i <= ticks.tickCount; i++) {
        const value = ticks.min + i * ticks.step;
        const y = this.valueToY(value, ticks.min, ticks.max, grid);
        if (splitLine.show !== false) {
          this.ctx.beginPath();
          this.ctx.setStrokeStyle(lineColor);
          this.ctx.setLineWidth(1);
          this.ctx.moveTo(grid.left, y);
          this.ctx.lineTo(this.canvasWidth - grid.right, y);
          this.ctx.stroke();
        }

        this.ctx.setFillStyle(labelColor);
        this.ctx.fillText(
          this.formatAxisLabel(value, axisLabel, digits, compact),
          grid.left - 8,
          y
        );
      }
    },
    drawXAxis(visibleCandles, grid, xAxisOption) {
      const axisLabel = xAxisOption?.axisLabel || {};
      const labelColor = axisLabel.color || '#666';
      const axisLineColor = xAxisOption?.axisLine?.lineStyle?.color || '#d1d5db';
      const y = grid.top + grid.height;
      this.ctx.beginPath();
      this.ctx.setStrokeStyle(axisLineColor);
      this.ctx.setLineWidth(1);
      this.ctx.moveTo(grid.left, y);
      this.ctx.lineTo(this.canvasWidth - grid.right, y);
      this.ctx.stroke();

      const count = visibleCandles.items.length;
      if (count <= 0) {
        return;
      }
      const labelCount = count <= 8 ? count : Math.min(6, count);
      const step = count <= 8 ? 1 : Math.max(1, Math.floor((count - 1) / Math.max(labelCount - 1, 1)));
      const indices = [];
      for (let i = 0; i < count; i += step) {
        indices.push(i);
      }
      if (indices[indices.length - 1] !== count - 1) {
        indices.push(count - 1);
      }

      this.ctx.setFontSize(axisLabel.fontSize || 11);
      this.ctx.setFillStyle(labelColor);
      this.ctx.setTextAlign('center');
      this.ctx.setTextBaseline('top');
      indices.forEach((index) => {
        const item = visibleCandles.items[index];
        if (!item) {
          return;
        }
        const x = this.getCategoryCenter(index, count, grid);
        this.ctx.fillText(this.truncateText(item.category, 10), x, y + 8);
      });
    },
    drawCandles(candlestickSeries, visibleCandles, priceRange, grid) {
      const itemStyle = candlestickSeries?.itemStyle || {};
      const upColor = itemStyle.color || '#ec0000';
      const upBorderColor = itemStyle.borderColor || '#8A0000';
      const downColor = itemStyle.color0 || '#00da3c';
      const downBorderColor = itemStyle.borderColor0 || '#008F28';
      const candleWidth = Math.max(3, Math.min(18, (grid.width / Math.max(visibleCandles.items.length, 1)) * 0.66));
      this.candleRegions = [];

      visibleCandles.items.forEach((item) => {
        const centerX = this.getCategoryCenter(item.localIndex, visibleCandles.items.length, grid);
        const openY = this.valueToY(item.open, priceRange.min, priceRange.max, grid);
        const closeY = this.valueToY(item.close, priceRange.min, priceRange.max, grid);
        const highY = this.valueToY(item.high, priceRange.min, priceRange.max, grid);
        const lowY = this.valueToY(item.low, priceRange.min, priceRange.max, grid);
        const isUp = item.close >= item.open;
        const fillColor = isUp ? upColor : downColor;
        const strokeColor = isUp ? upBorderColor : downBorderColor;
        const bodyTop = Math.min(openY, closeY);
        const bodyHeight = Math.max(1, Math.abs(closeY - openY));
        const bodyX = centerX - candleWidth / 2;
        const regionHeight = Math.max(6, lowY - highY);

        this.ctx.beginPath();
        this.ctx.setStrokeStyle(strokeColor);
        this.ctx.setLineWidth(1);
        this.ctx.moveTo(centerX, highY);
        this.ctx.lineTo(centerX, lowY);
        this.ctx.stroke();

        if (Math.abs(closeY - openY) < 1) {
          this.ctx.beginPath();
          this.ctx.moveTo(bodyX, bodyTop);
          this.ctx.lineTo(bodyX + candleWidth, bodyTop);
          this.ctx.stroke();
        } else {
          this.ctx.setFillStyle(fillColor);
          this.ctx.fillRect(bodyX, bodyTop, candleWidth, bodyHeight);
          this.ctx.setStrokeStyle(strokeColor);
          this.ctx.strokeRect(bodyX, bodyTop, candleWidth, bodyHeight);
        }

        this.candleRegions.push({
          x: bodyX,
          y: highY,
          width: candleWidth,
          height: regionHeight,
          centerX,
          centerY: bodyTop + bodyHeight / 2,
          item,
          isUp,
          color: fillColor
        });
      });
    },
    drawStraightSegment(points) {
      this.ctx.beginPath();
      points.forEach((point, index) => {
        if (index === 0) {
          this.ctx.moveTo(point.x, point.y);
        } else {
          this.ctx.lineTo(point.x, point.y);
        }
      });
    },
    drawSmoothSegment(points) {
      if (!points || points.length < 2) {
        return;
      }
      this.ctx.beginPath();
      this.ctx.moveTo(points[0].x, points[0].y);
      if (points.length === 2) {
        const controlX = (points[0].x + points[1].x) / 2;
        const controlY = (points[0].y + points[1].y) / 2;
        this.ctx.quadraticCurveTo(controlX, controlY, points[1].x, points[1].y);
        return;
      }
      for (let i = 0; i < points.length - 1; i++) {
        const p0 = i === 0 ? points[0] : points[i - 1];
        const p1 = points[i];
        const p2 = points[i + 1];
        const p3 = i === points.length - 2 ? points[i + 1] : points[i + 2];
        const cp1x = p1.x + (p2.x - p0.x) / 6;
        const cp1y = p1.y + (p2.y - p0.y) / 6;
        const cp2x = p2.x - (p3.x - p1.x) / 6;
        const cp2y = p2.y - (p3.y - p1.y) / 6;
        this.ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, p2.x, p2.y);
      }
    },
    drawOverlayLines(lineSeries, visibleCandles, priceRange, grid) {
      lineSeries.forEach((serie, seriesIndex) => {
        const data = Array.isArray(serie.data) ? serie.data : [];
        const segments = [];
        let currentSegment = [];

        visibleCandles.items.forEach((item) => {
          const value = this.parseNumber(data[item.index]);
          if (value === null) {
            if (currentSegment.length > 1) {
              segments.push(currentSegment);
            }
            currentSegment = [];
            return;
          }
          currentSegment.push({
            x: this.getCategoryCenter(item.localIndex, visibleCandles.items.length, grid),
            y: this.valueToY(value, priceRange.min, priceRange.max, grid),
            value
          });
        });
        if (currentSegment.length > 1) {
          segments.push(currentSegment);
        }

        const color = this.getSeriesColor(serie, seriesIndex);
        const opacity = serie.lineStyle?.opacity !== undefined ? serie.lineStyle.opacity : 0.9;
        this.ctx.setStrokeStyle(this.applyOpacity(color, opacity));
        this.ctx.setLineWidth(serie.lineStyle?.width || 1.5);
        this.ctx.setLineJoin('round');
        this.ctx.setLineCap('round');
        segments.forEach((segment) => {
          if (serie.smooth === true) {
            this.drawSmoothSegment(segment);
          } else {
            this.drawStraightSegment(segment);
          }
          this.ctx.stroke();
        });
      });
    },
    getBarColor(serie, normalizedItem, candleItem, seriesIndex) {
      if (normalizedItem.itemStyle?.color) {
        return normalizedItem.itemStyle.color;
      }
      if (serie.itemStyle?.color) {
        return serie.itemStyle.color;
      }
      if (candleItem) {
        return candleItem.close >= candleItem.open
          ? '#ec0000'
          : '#00da3c';
      }
      return chartHelper.getColor(seriesIndex + 1);
    },
    drawVolumeBars(barSeries, visibleCandles, volumeRange, grid) {
      this.volumeRegions = [];
      if (!grid || barSeries.length === 0) {
        return;
      }

      const count = Math.max(visibleCandles.items.length, 1);
      const categoryWidth = grid.width / count;
      const gap = 2;
      const groupWidth = Math.max(4, Math.min(22, categoryWidth * 0.72));
      const barWidth = Math.max(2, (groupWidth - gap * Math.max(barSeries.length - 1, 0)) / Math.max(barSeries.length, 1));
      const baselineY = this.valueToY(volumeRange.baseline, volumeRange.min, volumeRange.max, grid);

      visibleCandles.items.forEach((item) => {
        const centerX = this.getCategoryCenter(item.localIndex, count, grid);
        const totalWidth = barSeries.length * barWidth + gap * Math.max(barSeries.length - 1, 0);
        const startX = centerX - totalWidth / 2;
        const groupValues = [];

        barSeries.forEach((serie, seriesIndex) => {
          const data = Array.isArray(serie.data) ? serie.data : [];
          const normalizedItem = this.normalizeSeriesValue(data[item.index]);
          if (normalizedItem.value === null) {
            return;
          }
          const valueY = this.valueToY(normalizedItem.value, volumeRange.min, volumeRange.max, grid);
          const barX = startX + seriesIndex * (barWidth + gap);
          const barY = Math.min(valueY, baselineY);
          const barHeight = Math.max(1, Math.abs(baselineY - valueY));
          const barColor = this.getBarColor(serie, normalizedItem, item, seriesIndex);

          this.ctx.setFillStyle(this.applyOpacity(barColor, serie.itemStyle?.opacity !== undefined ? serie.itemStyle.opacity : 0.88));
          this.ctx.fillRect(barX, barY, barWidth, barHeight);

          groupValues.push({
            seriesName: serie.name || `bar-${seriesIndex + 1}`,
            value: normalizedItem.value,
            color: barColor,
            raw: normalizedItem.raw
          });
        });

        if (groupValues.length > 0) {
          this.volumeRegions.push({
            x: centerX - groupWidth / 2,
            y: grid.top,
            width: groupWidth,
            height: grid.height,
            centerX,
            item,
            values: groupValues
          });
        }
      });
    },
    resolveMarkTarget(config, visibleCandles, priceRange, grid) {
      if (!config) {
        return null;
      }

      if (Array.isArray(config.coord) && config.coord.length >= 2) {
        const category = config.coord[0];
        const value = Number(config.coord[1]);
        const target = visibleCandles.items.find((item) => String(item.category) === String(category));
        if (!target) {
          return null;
        }
        return {
          x: this.getCategoryCenter(target.localIndex, visibleCandles.items.length, grid),
          y: this.valueToY(value, priceRange.min, priceRange.max, grid),
          value,
          item: target,
          name: config.name || ''
        };
      }

      const valueDim = config.valueDim || 'close';
      if (config.type === 'average') {
        let sum = 0;
        let count = 0;
        visibleCandles.items.forEach((item) => {
          sum += this.getValueByDim(item, valueDim);
          count += 1;
        });
        const average = count > 0 ? sum / count : 0;
        const target = visibleCandles.items[Math.floor(visibleCandles.items.length / 2)];
        if (!target) {
          return null;
        }
        return {
          x: this.getCategoryCenter(target.localIndex, visibleCandles.items.length, grid),
          y: this.valueToY(average, priceRange.min, priceRange.max, grid),
          value: average,
          item: target,
          name: config.name || 'average'
        };
      }

      if (config.type === 'max' || config.type === 'min') {
        let target = null;
        visibleCandles.items.forEach((item) => {
          const currentValue = this.getValueByDim(item, valueDim);
          if (!target) {
            target = item;
            return;
          }
          if (config.type === 'max' && currentValue > this.getValueByDim(target, valueDim)) {
            target = item;
          }
          if (config.type === 'min' && currentValue < this.getValueByDim(target, valueDim)) {
            target = item;
          }
        });
        if (!target) {
          return null;
        }
        const targetValue = this.getValueByDim(target, valueDim);
        return {
          x: this.getCategoryCenter(target.localIndex, visibleCandles.items.length, grid),
          y: this.valueToY(targetValue, priceRange.min, priceRange.max, grid),
          value: targetValue,
          item: target,
          name: config.name || config.type
        };
      }

      return null;
    },
    formatMarkPointLabel(config, value) {
      const formatter = config?.label?.formatter;
      if (typeof formatter === 'function') {
        try {
          return String(formatter({ value }));
        } catch (error) {
          return `${Math.round(value)}`;
        }
      }
      return `${Math.round(value)}`;
    },
    drawMarkPoints(candlestickSeries, visibleCandles, priceRange, grid) {
      const markPoint = candlestickSeries?.markPoint;
      if (!markPoint || !Array.isArray(markPoint.data)) {
        this.markPointRegions = [];
        return;
      }
      this.markPointRegions = [];
      markPoint.data.forEach((config) => {
        const target = this.resolveMarkTarget(config, visibleCandles, priceRange, grid);
        if (!target) {
          return;
        }
        const pointColor = config.itemStyle?.color || '#293c55';
        const labelText = this.formatMarkPointLabel(config, target.value);

        this.ctx.beginPath();
        this.ctx.setFillStyle(pointColor);
        this.ctx.arc(target.x, target.y, 5, 0, Math.PI * 2);
        this.ctx.fill();

        this.ctx.setFontSize(10);
        this.ctx.setFillStyle(pointColor);
        this.ctx.setTextAlign('center');
        this.ctx.setTextBaseline('bottom');
        this.ctx.fillText(labelText, target.x, target.y - 8);

        this.markPointRegions.push({
          x: target.x - 10,
          y: target.y - 20,
          width: 20,
          height: 24,
          value: target.value,
          config,
          category: target.item?.category,
          centerX: target.x,
          centerY: target.y
        });
      });
    },
    drawMarkLines(candlestickSeries, visibleCandles, priceRange, grid) {
      const markLine = candlestickSeries?.markLine;
      if (!markLine || !Array.isArray(markLine.data)) {
        return;
      }
      markLine.data.forEach((config) => {
        if (Array.isArray(config) && config.length === 2) {
          const start = this.resolveMarkTarget(config[0], visibleCandles, priceRange, grid);
          const end = this.resolveMarkTarget(config[1], visibleCandles, priceRange, grid);
          if (!start || !end) {
            return;
          }
          this.ctx.beginPath();
          this.ctx.setStrokeStyle('#64748b');
          this.ctx.setLineWidth(1);
          this.ctx.moveTo(start.x, start.y);
          this.ctx.lineTo(end.x, end.y);
          this.ctx.stroke();
          return;
        }

        const target = this.resolveMarkTarget(config, visibleCandles, priceRange, grid);
        if (!target) {
          return;
        }
        this.ctx.beginPath();
        this.ctx.setStrokeStyle('#94a3b8');
        this.ctx.setLineWidth(1);
        this.ctx.moveTo(grid.left, target.y);
        this.ctx.lineTo(this.canvasWidth - grid.right, target.y);
        this.ctx.stroke();

        this.ctx.setFontSize(10);
        this.ctx.setFillStyle('#64748b');
        this.ctx.setTextAlign('left');
        this.ctx.setTextBaseline('bottom');
        this.ctx.fillText(config.name || `${Math.round(target.value)}`, grid.left + 4, target.y - 2);
      });
    },
    getSliderSelection() {
      if (!this.sliderRect) {
        return null;
      }
      const selectedX = this.sliderRect.x + this.sliderRect.width * this.zoomStart / 100;
      const selectedWidth = Math.max(8, this.sliderRect.width * (this.zoomEnd - this.zoomStart) / 100);
      return {
        selectedX,
        selectedWidth
      };
    },
    drawDataZoomSlider(showSlider) {
      if (!showSlider) {
        this.sliderRect = null;
        return;
      }
      const x = this.chartGrid.left;
      const y = this.canvasHeight - 22;
      const width = this.chartGrid.width;
      const height = 12;
      this.sliderRect = { x, y, width, height };

      this.ctx.setFillStyle('#e2e8f0');
      this.ctx.fillRect(x, y, width, height);

      const selection = this.getSliderSelection();
      if (!selection) {
        return;
      }
      this.ctx.setFillStyle('rgba(59, 130, 246, 0.28)');
      this.ctx.fillRect(selection.selectedX, y, selection.selectedWidth, height);

      this.ctx.setFillStyle('#3b82f6');
      this.ctx.fillRect(selection.selectedX, y - 2, 4, height + 4);
      this.ctx.fillRect(selection.selectedX + selection.selectedWidth - 4, y - 2, 4, height + 4);
    },
    isPointInRect(x, y, rect) {
      return rect && x >= rect.x && x <= rect.x + rect.width && y >= rect.y && y <= rect.y + rect.height;
    },
    isPointInGrid(x, y, grid) {
      return !!grid && x >= grid.left && x <= grid.left + grid.width && y >= grid.top && y <= grid.top + grid.height;
    },
    buildPointerPayload(item, touchX, touchY, chartState) {
      const { candlestickSeries, lineSeries, barSeries, priceGrid, volumeGrid, priceRange } = chartState;
      const color = item.close >= item.open
        ? (candlestickSeries.itemStyle?.color || '#ec0000')
        : (candlestickSeries.itemStyle?.color0 || '#00da3c');
      const centerX = this.getCategoryCenter(item.localIndex, chartState.visibleCandles.items.length, priceGrid);
      const priceY = this.isPointInGrid(touchX, touchY, priceGrid)
        ? Math.max(priceGrid.top, Math.min(priceGrid.top + priceGrid.height, touchY))
        : null;
      const lineValues = lineSeries
        .map((serie, index) => {
          const data = Array.isArray(serie.data) ? serie.data : [];
          const value = this.parseNumber(data[item.index]);
          if (value === null) {
            return null;
          }
          return {
            seriesType: 'line',
            seriesName: serie.name || `line-${index + 1}`,
            value,
            color: this.getSeriesColor(serie, index)
          };
        })
        .filter(Boolean);

      const barValues = barSeries
        .map((serie, index) => {
          const data = Array.isArray(serie.data) ? serie.data : [];
          const normalizedItem = this.normalizeSeriesValue(data[item.index]);
          if (normalizedItem.value === null) {
            return null;
          }
          return {
            seriesType: 'bar',
            seriesName: serie.name || `bar-${index + 1}`,
            value: normalizedItem.value,
            color: this.getBarColor(serie, normalizedItem, item, index),
            raw: normalizedItem.raw
          };
        })
        .filter(Boolean);

      const previousRaw = Array.isArray(candlestickSeries.data) && item.index > 0
        ? candlestickSeries.data[item.index - 1]
        : null;
      const previousClose = Array.isArray(previousRaw) && previousRaw.length > 1
        ? Number(previousRaw[1])
        : null;
      const change = previousClose !== null ? item.close - previousClose : item.close - item.open;
      const changePercent = previousClose
        ? (change / previousClose) * 100
        : (item.open ? ((item.close - item.open) / item.open) * 100 : 0);

      return {
        componentType: 'series',
        seriesType: 'candlestick',
        seriesName: candlestickSeries.name || 'candlestick',
        name: item.category,
        axisValue: item.category,
        axisValueLabel: String(item.category),
        dataIndex: item.index,
        value: item.raw,
        open: item.open,
        close: item.close,
        lowest: item.low,
        highest: item.high,
        color,
        previousClose,
        change,
        changePercent,
        lineValues,
        barValues,
        volume: barValues[0]?.value ?? null,
        x: centerX,
        priceY,
        pointerValue: priceY !== null ? this.yToValue(priceY, priceRange.min, priceRange.max, priceGrid) : null,
        event: {
          offsetX: centerX,
          offsetY: priceY !== null ? priceY : this.valueToY(item.close, priceRange.min, priceRange.max, priceGrid)
        },
        withinPrice: this.isPointInGrid(touchX, touchY, priceGrid),
        withinVolume: this.isPointInGrid(touchX, touchY, volumeGrid),
        item
      };
    },
    syncActivePointer(chartState) {
      if (!this.activePointer || !chartState) {
        return;
      }
      const matched = chartState.visibleCandles.items.find((item) => item.index === this.activePointer.dataIndex);
      if (!matched) {
        this.activePointer = null;
        return;
      }
      const nextTouchX = this.activePointer.x || this.touchInfo.lastX || 0;
      const nextTouchY = this.activePointer.priceY || this.touchInfo.lastY || 0;
      this.activePointer = this.buildPointerPayload(matched, nextTouchX, nextTouchY, chartState);
    },
    updateActivePointer(touchX, touchY, emitTooltip = true) {
      const chartState = this.renderState;
      if (!chartState || !chartState.visibleCandles?.items?.length) {
        return null;
      }
      const targetGrid = this.isPointInGrid(touchX, touchY, chartState.priceGrid) || this.isPointInGrid(touchX, touchY, chartState.volumeGrid)
        ? chartState.priceGrid
        : null;
      if (!targetGrid) {
        return null;
      }

      const count = chartState.visibleCandles.items.length;
      const localIndex = Math.max(
        0,
        Math.min(
          count - 1,
          Math.floor(((touchX - chartState.priceGrid.left) / Math.max(chartState.priceGrid.width, 1)) * count)
        )
      );
      const item = chartState.visibleCandles.items[localIndex];
      if (!item) {
        return null;
      }

      this.touchInfo.lastX = touchX;
      this.touchInfo.lastY = touchY;
      this.activePointer = this.buildPointerPayload(item, touchX, touchY, chartState);
      if (emitTooltip) {
        this.$emit('tooltipShow', this.activePointer);
      }
      this.drawChart();
      return this.activePointer;
    },
    drawCrosshair(chartState) {
      if (!this.activePointer || !chartState) {
        return;
      }
      const tooltipOption = this.option?.tooltip || {};
      const axisPointer = tooltipOption.axisPointer || {};
      if (axisPointer.show === false) {
        return;
      }

      const lineColor = axisPointer.lineStyle?.color || 'rgba(71, 85, 105, 0.75)';
      const lineWidth = axisPointer.lineStyle?.width || 1;
      const bottom = chartState.volumeGrid
        ? chartState.volumeGrid.top + chartState.volumeGrid.height
        : chartState.priceGrid.top + chartState.priceGrid.height;

      this.ctx.beginPath();
      this.ctx.setStrokeStyle(lineColor);
      this.ctx.setLineWidth(lineWidth);
      this.ctx.moveTo(this.activePointer.x, chartState.priceGrid.top);
      this.ctx.lineTo(this.activePointer.x, bottom);
      this.ctx.stroke();

      if ((axisPointer.type === 'cross' || axisPointer.type === undefined) && this.activePointer.priceY !== null) {
        this.ctx.beginPath();
        this.ctx.setStrokeStyle(this.applyOpacity(lineColor, 0.8));
        this.ctx.setLineWidth(lineWidth);
        this.ctx.moveTo(chartState.priceGrid.left, this.activePointer.priceY);
        this.ctx.lineTo(this.canvasWidth - chartState.priceGrid.right, this.activePointer.priceY);
        this.ctx.stroke();
      }
    },
    getTooltipLines(pointer) {
      const changeSign = pointer.change > 0 ? '+' : '';
      const lines = [
        { text: String(pointer.name), color: '#f8fafc', fontSize: 12, weight: 'bold' },
        { text: `开 ${this.formatNumber(pointer.open)}  收 ${this.formatNumber(pointer.close)}`, color: '#e2e8f0', fontSize: 11 },
        { text: `低 ${this.formatNumber(pointer.lowest)}  高 ${this.formatNumber(pointer.highest)}`, color: '#cbd5e1', fontSize: 11 },
        {
          text: `涨跌 ${changeSign}${this.formatNumber(pointer.change)} (${changeSign}${this.formatNumber(pointer.changePercent)}%)`,
          color: pointer.change >= 0 ? '#f87171' : '#4ade80',
          fontSize: 11
        }
      ];

      if (pointer.volume !== null) {
        lines.push({
          text: `成交量 ${this.formatCompactValue(pointer.volume, 2)}`,
          color: '#93c5fd',
          fontSize: 11
        });
      }

      pointer.lineValues.forEach((lineItem) => {
        lines.push({
          text: `${lineItem.seriesName} ${this.formatNumber(lineItem.value)}`,
          color: lineItem.color,
          fontSize: 11
        });
      });
      return lines;
    },
    drawTooltipBox(chartState) {
      if (!this.activePointer || !chartState) {
        return;
      }
      const tooltipOption = this.option?.tooltip || {};
      if (tooltipOption.show === false || tooltipOption.showContent === false) {
        return;
      }

      const lines = this.getTooltipLines(this.activePointer);
      const paddingX = 10;
      const paddingY = 8;
      const lineGap = 6;
      let boxWidth = 0;
      let boxHeight = paddingY * 2 - lineGap;
      lines.forEach((line) => {
        boxWidth = Math.max(boxWidth, this.measureTextWidth(line.text, line.fontSize));
        boxHeight += line.fontSize + lineGap;
      });
      boxWidth += paddingX * 2;

      let boxX = this.activePointer.x + 12;
      if (boxX + boxWidth > this.canvasWidth - 8) {
        boxX = this.activePointer.x - boxWidth - 12;
      }
      boxX = Math.max(8, boxX);

      let boxY = chartState.priceGrid.top + 10;
      const maxBoxY = Math.max(chartState.priceGrid.top + 4, chartState.priceGrid.top + chartState.priceGrid.height - boxHeight - 4);
      boxY = Math.min(boxY, maxBoxY);

      this.ctx.setFillStyle('rgba(15, 23, 42, 0.88)');
      this.ctx.fillRect(boxX, boxY, boxWidth, boxHeight);
      this.ctx.setStrokeStyle('rgba(148, 163, 184, 0.5)');
      this.ctx.setLineWidth(1);
      this.ctx.strokeRect(boxX, boxY, boxWidth, boxHeight);

      let currentY = boxY + paddingY;
      lines.forEach((line) => {
        this.ctx.setFontSize(line.fontSize);
        this.ctx.setFillStyle(line.color);
        this.ctx.setTextAlign('left');
        this.ctx.setTextBaseline('top');
        this.ctx.fillText(line.text, boxX + paddingX, currentY);
        currentY += line.fontSize + lineGap;
      });
    },
    drawChart() {
      if (!this.ctx || !this.option) {
        return;
      }

      const { candlestickSeries, lineSeries, barSeries, allSeries } = this.getSeriesCollection();
      if (!candlestickSeries) {
        this.ctx.clearRect(0, 0, this.canvasWidth, this.canvasHeight);
        this.renderState = null;
        return;
      }

      const xAxisSource = this.option?.xAxis || {};
      const yAxisSource = this.option?.yAxis || {};
      const categories = Array.isArray(this.getAxisOption(xAxisSource, 0).data)
        ? this.getAxisOption(xAxisSource, 0).data
        : (Array.isArray(candlestickSeries.data) ? candlestickSeries.data.map((_, index) => `${index}`) : []);
      const showSlider = Array.isArray(this.option?.dataZoom)
        ? this.option.dataZoom.some((item) => item.type === 'slider' && item.show !== false)
        : false;

      this.ctx.clearRect(0, 0, this.canvasWidth, this.canvasHeight);
      if (this.option.backgroundColor) {
        this.ctx.setFillStyle(this.option.backgroundColor);
        this.ctx.fillRect(0, 0, this.canvasWidth, this.canvasHeight);
      }

      const titleHeight = this.drawTitle();
      const baseGrid = this.createBaseGrid(titleHeight, showSlider);

      const legend = this.option?.legend;
      const legendData = Array.isArray(legend?.data) && legend.data.length > 0
        ? legend.data
        : allSeries.map((serie) => serie.name).filter(Boolean);
      if (legend && legend.show !== false && legendData.length > 0) {
        chartHelper.drawLegend(
          this.ctx,
          { ...legend, data: legendData },
          baseGrid,
          this.canvasWidth,
          this.getLegendColors(allSeries, legendData),
          this.canvasHeight,
          titleHeight
        );
      }

      this.chartGrid = {
        ...baseGrid,
        width: Math.max(0, this.canvasWidth - baseGrid.left - baseGrid.right),
        height: Math.max(0, this.canvasHeight - baseGrid.top - baseGrid.bottom)
      };

      const hasVolume = barSeries.length > 0;
      const { priceGrid, volumeGrid } = this.createPanelLayout(this.chartGrid, hasVolume);
      this.priceGrid = priceGrid;
      this.volumeGrid = volumeGrid;

      const visibleCandles = this.buildVisibleCandles(candlestickSeries, categories);
      if (visibleCandles.items.length === 0) {
        this.drawDataZoomSlider(showSlider);
        this.ctx.draw();
        this.renderState = null;
        return;
      }

      const priceYAxis = this.getAxisOption(yAxisSource, 0);
      const volumeYAxis = this.getAxisOption(yAxisSource, 1);
      const priceXAxis = this.getAxisOption(xAxisSource, 0);
      const volumeXAxis = this.getAxisOption(xAxisSource, 1);
      const priceRange = this.calculatePriceRange(visibleCandles, lineSeries, priceYAxis);
      const priceTicks = chartHelper.calculateYAxisTicks(priceRange.min, priceRange.max, 5);
      priceRange.min = priceTicks.min;
      priceRange.max = priceTicks.max;

      let volumeRange = null;
      let volumeTicks = null;
      if (hasVolume && volumeGrid) {
        volumeRange = this.calculateVolumeRange(visibleCandles, barSeries, volumeYAxis);
        volumeTicks = chartHelper.calculateYAxisTicks(volumeRange.min, volumeRange.max, 2);
        volumeRange.min = volumeTicks.min;
        volumeRange.max = volumeTicks.max;
        volumeRange.baseline = Math.min(Math.max(0, volumeRange.min), volumeRange.max);
      }

      this.renderState = {
        candlestickSeries,
        lineSeries,
        barSeries,
        visibleCandles,
        priceGrid,
        volumeGrid,
        priceRange,
        volumeRange,
        priceTicks,
        volumeTicks
      };
      this.syncActivePointer(this.renderState);

      this.drawYAxisSplitArea(priceTicks, priceGrid, priceYAxis);
      this.drawYAxis(priceTicks, priceGrid, priceYAxis, false);
      this.drawCandles(candlestickSeries, visibleCandles, priceRange, priceGrid);
      this.drawOverlayLines(lineSeries, visibleCandles, priceRange, priceGrid);
      this.drawMarkLines(candlestickSeries, visibleCandles, priceRange, priceGrid);
      this.drawMarkPoints(candlestickSeries, visibleCandles, priceRange, priceGrid);

      if (hasVolume && volumeGrid && volumeRange && volumeTicks) {
        this.drawYAxis(volumeTicks, volumeGrid, volumeYAxis, true);
        this.drawVolumeBars(barSeries, visibleCandles, volumeRange, volumeGrid);
        this.drawXAxis(visibleCandles, volumeGrid, volumeXAxis?.data ? volumeXAxis : priceXAxis);

        this.ctx.beginPath();
        this.ctx.setStrokeStyle('#cbd5e1');
        this.ctx.setLineWidth(1);
        this.ctx.moveTo(priceGrid.left, priceGrid.top + priceGrid.height);
        this.ctx.lineTo(this.canvasWidth - priceGrid.right, priceGrid.top + priceGrid.height);
        this.ctx.stroke();
      } else {
        this.drawXAxis(visibleCandles, priceGrid, priceXAxis);
      }

      this.drawCrosshair(this.renderState);
      this.drawTooltipBox(this.renderState);
      this.drawDataZoomSlider(showSlider);
      this.ctx.draw();
    },
    setZoom(start, end) {
      const safeStart = Math.max(0, Math.min(99, Number(start)));
      const safeEnd = Math.max(safeStart + 1, Math.min(100, Number(end)));
      this.zoomStart = safeStart;
      this.zoomEnd = safeEnd;
      this.$emit('zoom', { start: this.zoomStart, end: this.zoomEnd });
      this.drawChart();
    },
    setOption(option, notMerge = false) {
      if (notMerge) {
        this.$emit('update:option', option);
        this.drawChart();
        return;
      }
      try {
        const newOption = JSON.parse(JSON.stringify(this.option || {}));
        Object.assign(newOption, option);
        this.$emit('update:option', newOption);
      } catch (error) {
        console.error('合并 K 线配置失败:', error);
      }
      this.drawChart();
    },
    resize() {
      this.initCanvas();
    },
    resolveSliderDragMode(x) {
      const selection = this.getSliderSelection();
      if (!selection) {
        return 'move';
      }
      const handleArea = 12;
      if (Math.abs(x - selection.selectedX) <= handleArea) {
        return 'start';
      }
      if (Math.abs(x - (selection.selectedX + selection.selectedWidth)) <= handleArea) {
        return 'end';
      }
      if (x >= selection.selectedX && x <= selection.selectedX + selection.selectedWidth) {
        return 'move';
      }
      return 'jump';
    },
    applySliderDrag(x) {
      if (!this.sliderRect) {
        return;
      }
      const mode = this.touchInfo.sliderDragMode || 'move';
      const deltaPercent = ((x || 0) - this.touchInfo.startX) / Math.max(this.sliderRect.width, 1) * 100;
      const minSpan = Math.max(1, Number(this.option?.dataZoom?.[0]?.minSpan || 5));
      const zoomSpan = this.touchInfo.originZoomEnd - this.touchInfo.originZoomStart;

      let nextStart = this.zoomStart;
      let nextEnd = this.zoomEnd;

      if (mode === 'start') {
        nextStart = Math.max(0, Math.min(this.touchInfo.originZoomEnd - minSpan, this.touchInfo.originZoomStart + deltaPercent));
        nextEnd = this.touchInfo.originZoomEnd;
      } else if (mode === 'end') {
        nextStart = this.touchInfo.originZoomStart;
        nextEnd = Math.min(100, Math.max(this.touchInfo.originZoomStart + minSpan, this.touchInfo.originZoomEnd + deltaPercent));
      } else if (mode === 'jump') {
        const centerPercent = ((x - this.sliderRect.x) / Math.max(this.sliderRect.width, 1)) * 100;
        nextStart = Math.max(0, Math.min(100 - zoomSpan, centerPercent - zoomSpan / 2));
        nextEnd = nextStart + zoomSpan;
      } else {
        nextStart = Math.max(0, Math.min(100 - zoomSpan, this.touchInfo.originZoomStart + deltaPercent));
        nextEnd = nextStart + zoomSpan;
      }

      this.zoomStart = nextStart;
      this.zoomEnd = nextEnd;
      this.drawChart();
    },
    buildClickPayload(pointer, touchX, touchY) {
      if (!pointer) {
        return null;
      }
      return {
        componentType: 'series',
        seriesType: 'candlestick',
        seriesName: pointer.seriesName,
        name: pointer.name,
        dataIndex: pointer.dataIndex,
        value: pointer.value,
        open: pointer.open,
        close: pointer.close,
        lowest: pointer.lowest,
        highest: pointer.highest,
        color: pointer.color,
        volume: pointer.volume,
        lineValues: pointer.lineValues,
        barValues: pointer.barValues,
        event: {
          offsetX: touchX,
          offsetY: touchY
        }
      };
    },
    handleTouchStart(event) {
      const touch = event.touches && event.touches.length > 0 ? event.touches[0] : null;
      if (!touch) {
        return;
      }
      this.touchInfo.startX = touch.x || 0;
      this.touchInfo.startY = touch.y || 0;
      this.touchInfo.lastX = touch.x || 0;
      this.touchInfo.lastY = touch.y || 0;
      this.touchInfo.originZoomStart = this.zoomStart;
      this.touchInfo.originZoomEnd = this.zoomEnd;
      this.touchInfo.draggingSlider = this.isPointInRect(this.touchInfo.startX, this.touchInfo.startY, this.sliderRect);
      this.touchInfo.sliderDragMode = this.touchInfo.draggingSlider
        ? this.resolveSliderDragMode(this.touchInfo.startX)
        : '';

      if (!this.touchInfo.draggingSlider) {
        this.updateActivePointer(this.touchInfo.startX, this.touchInfo.startY, true);
      }
    },
    handleTouchMove(event) {
      const touch = event.touches && event.touches.length > 0 ? event.touches[0] : null;
      if (!touch) {
        return;
      }

      const touchX = touch.x || 0;
      const touchY = touch.y || 0;
      if (this.touchInfo.draggingSlider && this.sliderRect) {
        this.applySliderDrag(touchX);
        return;
      }

      this.updateActivePointer(touchX, touchY, true);
      event.preventDefault && event.preventDefault();
    },
    handleTouchEnd(event) {
      const touch = event.changedTouches && event.changedTouches.length > 0 ? event.changedTouches[0] : null;
      if (!touch) {
        return;
      }
      const endX = touch.x || 0;
      const endY = touch.y || 0;

      if (this.touchInfo.draggingSlider) {
        this.touchInfo.draggingSlider = false;
        this.touchInfo.sliderDragMode = '';
        this.$emit('zoom', { start: this.zoomStart, end: this.zoomEnd });
        return;
      }

      const markPoint = this.markPointRegions.find((item) => this.isPointInRect(endX, endY, item));
      if (markPoint && Math.abs(endX - this.touchInfo.startX) <= 8 && Math.abs(endY - this.touchInfo.startY) <= 8) {
        const payload = {
          componentType: 'markPoint',
          name: markPoint.config?.name || '',
          value: markPoint.value,
          category: markPoint.category,
          event: { offsetX: markPoint.centerX, offsetY: markPoint.centerY }
        };
        this.$emit('click', payload);
        this.$emit('tooltipShow', payload);
        return;
      }

      const pointer = this.updateActivePointer(endX, endY, true) || this.activePointer;
      if (Math.abs(endX - this.touchInfo.startX) > 8 || Math.abs(endY - this.touchInfo.startY) > 8) {
        return;
      }

      const payload = this.buildClickPayload(pointer, endX, endY);
      if (payload) {
        this.$emit('click', payload);
        this.$emit('tooltipShow', payload);
      }
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initCanvas();
    });
  }
};
</script>

<style scoped>
.ly-charts-candlestick {
  position: relative;
}

.chart-canvas {
  width: 100%;
  height: 100%;
  display: block;
}
</style>
