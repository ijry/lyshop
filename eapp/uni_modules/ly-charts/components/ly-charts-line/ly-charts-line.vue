<!-- components/ly-charts-line/ly-charts-line.vue -->
<template>
  <view class="ly-charts-line" :style="{ width: containerWidth, height: containerHeight }">
    <canvas 
      class="chart-canvas" 
      :id="canvasId" 
      :canvas-id="canvasId"
      @touchstart="handleTouchStart"
      @touchmove="handleTouchMove"
      @touchend="handleTouchEnd"
    ></canvas>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { nextTick, getCurrentInstance } from 'vue';
import chartHelper from '../../libs/util/chartHelper.js';
const instance = getCurrentInstance()

// 定义props
const props = defineProps({
  // ECharts 风格的配置项
  option: {
    type: Object,
    default: () => ({})
  },
  // 容器高度
  height: {
    type: [String, Number],
    default: 400
  },
  // 容器宽度
  width: {
    type: [String, Number],
    default: '100%'
  }
});

// 响应式数据
const canvasId = ref('line-chart-' + Date.now());
const ctx = ref(null);
const canvasWidth = ref(0);
const canvasHeight = ref(0);
const grid = ref({ top: 10, right: 20, bottom: 25, left: 50 });
// 存储系列数据用于事件处理
const seriesData = ref([]);
// 触摸相关信息
const touchInfo = ref({
  startX: 0,
  startY: 0
});

// 计算属性
const containerHeight = computed(() => {
  return typeof props.height === 'number' ? props.height + 'px' : props.height;
});
const containerWidth = computed(() => {
  return typeof props.width === 'number' ? props.width + 'px' : props.width;
});

// 方法定义
const initCanvas = () => {
  try {
    const query = uni.createSelectorQuery().in(instance);
    query.select('#' + canvasId.value).boundingClientRect((res) => {
      if (res) {
        canvasWidth.value = res.width;
        canvasHeight.value = res.height;
        
        // 创建canvas上下文
        ctx.value = uni.createCanvasContext(canvasId.value, instance);
        if (!ctx.value) {
          console.error('无法获取canvas绘图上下文');
          return;
        }
        
        drawChart(props.option);
      } else {
        console.error('无法获取canvas信息');
      }
    }).exec();
  } catch (error) {
    console.error('初始化canvas失败:', error);
  }
};

const drawChart = (option) => {
  if (!ctx.value || !option) return;
  
  try {
    // 清空画布
    ctx.value.clearRect(0, 0, canvasWidth.value, canvasHeight.value);
    
    // 设置背景色
    if (option.backgroundColor) {
      ctx.value.setFillStyle(option.backgroundColor);
      ctx.value.fillRect(0, 0, canvasWidth.value, canvasHeight.value);
    }
    
    // 绘制标题
    let titleHeight = 0;
    if (option.title && option.title.text) {
      titleHeight = chartHelper.drawTitle(ctx.value, option.title, canvasWidth.value);
    }
    
    // 如果有标题，调整网格顶部边距
    if (titleHeight > 0) {
      grid.value.top = Math.max(grid.value.top, titleHeight + 10);
    }
    
    // 提取系列数据
    const series = option.series || [];
    const xAxis = option.xAxis || {};
    const yAxis = option.yAxis || {};
    
    // 处理x轴数据
    const xAxisData = chartHelper.processXAxisData(xAxis, series);
    
    // 计算数据范围
    const { minY, maxY } = chartHelper.calculateDataRange(series);
    
    // 获取X轴padding配置，默认为0
    const xAxisPadding = option.xAxisPadding || 0;

    // 绘制图例（在网格调整之前）
    let legendHeight = 0;
    if (option.legend && option.legend.data) {
      const legendOption = { 
        ...option.legend, 
      };
      legendHeight = chartHelper.drawLegend(
        ctx.value, 
        legendOption, 
        grid.value, 
        canvasWidth.value, 
        chartHelper.defaultColors,
        canvasHeight.value,
        titleHeight
      );
    }
    
    // 检查X轴标签是否需要旋转，如果需要则增加底部边距
    let needRotateLabels = false;
    if (xAxisData && xAxisData.length > 0) {
      const xAxisLabelColor = (xAxis && xAxis.axisLabel && xAxis.axisLabel.color) || '#666';
      const xAxisFontSize = (xAxis && xAxis.axisLabel && xAxis.axisLabel.fontSize) || 12;
      
      // 估算标签长度和可用宽度
      const chartWidth = canvasWidth.value - (grid.value.left || 0) - (grid.value.right || 0);
      const avgLabelWidth = chartWidth / xAxisData.length;
      const maxLabelLength = Math.floor(avgLabelWidth / (xAxisFontSize * 0.6));
      
      // 检查是否存在需要旋转的长标签
      for (const label of xAxisData) {
        if (String(label).length > maxLabelLength && maxLabelLength > 0) {
          needRotateLabels = true;
          break;
        }
      }
    }
    
    // 如果需要旋转标签，增加底部边距
    if (needRotateLabels) {
      grid.value.bottom = Math.max(grid.value.bottom, 60); // 增加底部边距以适应旋转标签
    }
    
    // 绘制网格
    const drawVerticalLines = false
    if (option.grid !== false) {
      chartHelper.drawGrid(ctx.value, grid.value, canvasWidth.value, canvasHeight.value, xAxisData.length, minY, maxY, drawVerticalLines, xAxisPadding);
    }

    // 绘制坐标轴
    chartHelper.drawAxis(ctx.value, grid.value, canvasWidth.value, canvasHeight.value, xAxisData, minY, maxY, xAxis, yAxis, 'line', xAxisPadding);
    
    // 绘制折线
    drawSeries(series, xAxisData, minY, maxY, chartHelper.adjustedYMin, chartHelper.adjustedYMax, xAxisPadding);
    
    // 绘制到画布
    ctx.value.draw();
  } catch (error) {
    console.error('绘制图表失败:', error);
  }
};

const drawSeries = (series, xAxisData, minY, maxY, adjustedYMin, adjustedYMax, xAxisPadding = 0) => {
  if (!series || series.length === 0) return;
  
  const chartWidth = canvasWidth.value - grid.value.left - grid.value.right;
  
  // 应用X轴padding
  const paddedChartWidth = chartWidth - 2 * xAxisPadding;
  
  const chartHeight = canvasHeight.value - grid.value.top - grid.value.bottom;
  
  seriesData.value = []; // 重置系列数据
  
  // 使用调整后的Y轴范围
  const useAdjustedY = adjustedYMin !== undefined && adjustedYMax !== undefined;
  const actualMinY = useAdjustedY ? adjustedYMin : minY;
  const actualMaxY = useAdjustedY ? adjustedYMax : maxY;
  
  // 首先绘制所有系列的填充区域（如果需要）
  series.forEach((serie, index) => {
    if (serie.type !== 'line') return;
    
    const color = serie.color || serie.itemStyle?.color || chartHelper.getColor(index);
    const fillArea = serie?.areaStyle || serie.fillArea;
    const smooth = serie.smooth === true;
    
    // 转换数据点为坐标
    const points = [];
    if (serie.data && Array.isArray(serie.data)) {
      serie.data.forEach((value, i) => {
        const actualValue = typeof value === 'object' && value !== null ? value.value : value;
        // 确保不会除以零，并应用X轴padding
        const x = grid.value.left + xAxisPadding + (xAxisData.length > 1 ? (i / (xAxisData.length - 1)) * paddedChartWidth : 0);
        // 使用调整后的Y轴范围计算Y坐标
        const y = grid.value.top + chartHeight - ((actualValue - actualMinY) / (actualMaxY - actualMinY ?? 1)) * chartHeight;
        points.push({ 
          x, 
          y, 
          value: actualValue, 
          name: xAxisData[i] || i,
          seriesName: serie.name ?? `Series ${index}`
        });
      });
    }
    
    // 绘制填充区域
    if (fillArea && points.length > 0) {
      drawFillArea(points, color, smooth, grid.value.top + chartHeight);
    }
  });
  
  // 然后绘制所有系列的线条和点
  series.forEach((serie, index) => {
    if (serie.type !== 'line') return;
    
    const color = serie.color || serie.itemStyle?.color || chartHelper.getColor(index);
    const showSymbol = serie.showSymbol !== false;
    const smooth = serie.smooth === true;
    const lineWidth = serie.lineStyle?.width || 2;
    // 获取标签配置
    const label = serie.label || {};
    const showLabel = label.show === true;
    const labelColor = label.color || '#666';
    const labelPosition = label.position || 'top';
    
    // 转换数据点为坐标
    const points = [];
    if (serie.data && Array.isArray(serie.data)) {
      serie.data.forEach((value, i) => {
        const actualValue = typeof value === 'object' && value !== null ? value.value : value;
        // 确保不会除以零，并应用X轴padding
        const x = grid.value.left + xAxisPadding + (xAxisData.length > 1 ? (i / (xAxisData.length - 1)) * paddedChartWidth : 0);
        // 使用调整后的Y轴范围计算Y坐标
        const y = grid.value.top + chartHeight - ((actualValue - actualMinY) / (actualMaxY - actualMinY || 1)) * chartHeight;
        points.push({ 
          x, 
          y, 
          value: actualValue, 
          name: xAxisData[i] || i,
          seriesName: serie.name || `Series ${index}`
        });
      });
    }
    
    // 保存系列数据用于事件处理
    seriesData.value.push({
      name: serie.name || `Series ${index}`,
      points,
      color
    });
    
    // 绘制线条
    if (points.length > 0) {
      ctx.value.setLineWidth(lineWidth);
      ctx.value.setStrokeStyle(color);
      ctx.value.setLineJoin('round');
      ctx.value.setLineCap('round');
      
      if (smooth && points.length > 1) {
        // 绘制平滑曲线
        drawSmoothLine(points);
      } else {
        // 绘制直线
        drawStraightLine(points);
      }
      
      ctx.value.stroke();
      
      // 绘制数据点
      if (showSymbol) {
        ctx.value.setFillStyle(color);
        points.forEach(point => {
          ctx.value.beginPath();
          ctx.value.arc(point.x, point.y, 4, 0, 2 * Math.PI);
          ctx.value.fill();
          
          ctx.value.beginPath();
          ctx.value.setFillStyle('#ffffff');
          ctx.value.arc(point.x, point.y, 2, 0, 2 * Math.PI);
          ctx.value.fill();
          
          ctx.value.setFillStyle(color);
        });
      }
      
      // 绘制标签
      if (showLabel) {
        drawLabels(points, labelColor, labelPosition);
      }
    }
  });
};

const drawFillArea = (points, color, smooth, chartBottom) => {
  if (!points || points.length === 0) return;
  
  ctx.value.beginPath();
  
  // 移动到起始点
  ctx.value.moveTo(points[0].x, chartBottom);
  ctx.value.lineTo(points[0].x, points[0].y);
  
  if (smooth && points.length > 1) {
    // 平滑填充区域
    for (let i = 0; i < points.length - 1; i++) {
      const p0 = i === 0 ? points[0] : points[i - 1];
      const p1 = points[i];
      const p2 = points[i + 1];
      const p3 = i === points.length - 2 ? points[i + 1] : points[i + 2];
      
      // 计算控制点
      const cp1x = p1.x + (p2.x - p0.x) / 6;
      const cp1y = p1.y + (p2.y - p0.y) / 6;
      const cp2x = p2.x - (p3.x - p1.x) / 6;
      const cp2y = p2.y - (p3.y - p1.y) / 6;
      
      ctx.value.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, p2.x, p2.y);
    }
  } else {
    // 直线填充区域
    for (let i = 1; i < points.length; i++) {
      ctx.value.lineTo(points[i].x, points[i].y);
    }
  }
  
  // 闭合路径到结束点
  ctx.value.lineTo(points[points.length - 1].x, chartBottom);
  ctx.value.closePath();
  
  // 设置填充样式
  if (typeof color === 'string' && color.startsWith('#')) {
    // 如果是十六进制颜色，转换为带透明度的颜色
    const hex = color.replace('#', '');
    const r = parseInt(hex.substring(0, 2), 16);
    const g = parseInt(hex.substring(2, 4), 16);
    const b = parseInt(hex.substring(4, 6), 16);
    ctx.value.setFillStyle(`rgba(${r}, ${g}, ${b}, 0.3)`);
  } else if (typeof color === 'string' && color.startsWith('rgb')) {
    // 如果是rgb颜色，转换为rgba颜色
    const rgb = color.replace('rgb(', '').replace(')', '');
    ctx.value.setFillStyle(`rgba(${rgb}, 0.3)`);
  } else {
    // 其他情况使用默认透明度
    ctx.value.setFillStyle('rgba(84, 112, 198, 0.3)');
  }
  
  ctx.value.fill();
};

const drawStraightLine = (points) => {
  if (!points || points.length === 0) return;
  
  ctx.value.beginPath();
  points.forEach((point, i) => {
    if (i === 0) {
      ctx.value.moveTo(point.x, point.y);
    } else {
      ctx.value.lineTo(point.x, point.y);
    }
  });
};

const drawSmoothLine = (points) => {
  if (!points || points.length < 2) return;
  
  ctx.value.beginPath();
  ctx.value.moveTo(points[0].x, points[0].y);
  
  // 对于只有两个点的情况，使用简单的二次贝塞尔曲线
  if (points.length === 2) {
    const controlX = (points[0].x + points[1].x) / 2;
    const controlY = (points[0].y + points[1].y) / 2;
    ctx.value.quadraticCurveTo(controlX, controlY, points[1].x, points[1].y);
    return;
  }
  
  // 对于多个点的情况，使用 Catmull-Rom 样条转贝塞尔曲线的方法
  // 这样可以确保所有点都在曲线上
  for (let i = 0; i < points.length - 1; i++) {
    const p0 = i === 0 ? points[0] : points[i - 1];
    const p1 = points[i];
    const p2 = points[i + 1];
    const p3 = i === points.length - 2 ? points[i + 1] : points[i + 2];
    
    // 计算控制点 (使用 Catmull-Rom 转贝塞尔)
    const cp1x = p1.x + (p2.x - p0.x) / 6;
    const cp1y = p1.y + (p2.y - p0.y) / 6;
    const cp2x = p2.x - (p3.x - p1.x) / 6;
    const cp2y = p2.y - (p3.y - p1.y) / 6;
    
    ctx.value.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, p2.x, p2.y);
  }
};

// 绘制标签
const drawLabels = (points, color, position) => {
  if (!points || points.length === 0) return;
  
  points.forEach(point => {
    const label = point.value.toString();
    ctx.value.setFontSize(10);
    ctx.value.setFillStyle(color);
    ctx.value.setTextAlign('center');
    ctx.value.setTextBaseline('bottom');
    
    let x = point.x;
    let y = point.y;
    
    // 根据位置调整标签坐标
    switch (position) {
      case 'top':
        y -= 6;
        break;
      case 'bottom':
        y += 6;
        ctx.value.setTextBaseline('top');
        break;
      case 'left':
        x -= 6;
        ctx.value.setTextAlign('right');
        ctx.value.setTextBaseline('middle');
        break;
      case 'right':
        x += 6;
        ctx.value.setTextAlign('left');
        ctx.value.setTextBaseline('middle');
        break;
      default:
        y -= 10;
    }
    
    ctx.value.fillText(label, x, y);
  });
};

// 触摸事件处理
const handleTouchStart = (e) => {
  if (e.touches && e.touches.length > 0) {
    const touch = e.touches[0];
    touchInfo.value.startX = touch.x || 0;
    touchInfo.value.startY = touch.y || 0;
  }
};

const handleTouchMove = (e) => {
  // 阻止默认行为，避免页面滚动
  e.preventDefault && e.preventDefault();
};

const handleTouchEnd = (e) => {
  if (e.changedTouches && e.changedTouches.length > 0) {
    const touch = e.changedTouches[0];
    const endX = touch.x || 0;
    const endY = touch.y || 0;
    
    // 判断是否为点击事件
    const deltaX = Math.abs(endX - touchInfo.value.startX);
    const deltaY = Math.abs(endY - touchInfo.value.startY);
    
    if (deltaX < 5 && deltaY < 5) {
      handleChartClick(endX, endY);
    }
  }
};

const handleChartClick = (x, y) => {
  // 查找最近的数据点
  let minDistance = Infinity;
  let closestPoint = null;
  
  if (seriesData.value && seriesData.value.length > 0) {
    seriesData.value.forEach(series => {
      if (series.points && series.points.length > 0) {
        series.points.forEach(point => {
          const distance = Math.sqrt(Math.pow(point.x - x, 2) + Math.pow(point.y - y, 2));
          if (distance < minDistance && distance < 20) { // 20px 的点击范围
            minDistance = distance;
            closestPoint = {
              seriesName: series.name,
              name: point.name,
              value: point.value,
              color: series.color,
              x: point.x,
              y: point.y
            };
          }
        });
      }
    });
  }
  
  if (closestPoint) {
    // 触发点击事件
    emit('click', {
      componentType: 'series',
      seriesType: 'line',
      seriesName: closestPoint.seriesName,
      name: closestPoint.name,
      value: closestPoint.value,
      color: closestPoint.color,
      event: {
        offsetX: closestPoint.x,
        offsetY: closestPoint.y
      }
    });
  }
};

// 定义emit
const emit = defineEmits(['click']);

// 提供类似 ECharts 的 setOption 方法
const setOption = (option, notMerge = false) => {
  if (notMerge) {
    drawChart(option);
  } else {
    // 简单合并配置
    try {
      const newOption = JSON.parse(JSON.stringify(props.option));
      Object.assign(newOption, option);
      drawChart(newOption);
    } catch (error) {
      console.error('合并配置失败:', error);
      drawChart(option);
    }
  }
};

// 提供类似 ECharts 的 resize 方法
const resize = () => {
  initCanvas();
};

// 监听option变化
watch(() => props.option, (newOption) => {
  drawChart(newOption);
}, { deep: true });

// 组件挂载后初始化
onMounted(() => {
  nextTick(() => {
    initCanvas();
  });
});

// 导出需要在模板中使用的变量和方法
defineExpose({
  setOption,
  resize
});
</script>

<style scoped>
.ly-charts-line {
  position: relative;
}

.chart-canvas {
  width: 100%;
  height: 100%;
  display: block;
}
</style>