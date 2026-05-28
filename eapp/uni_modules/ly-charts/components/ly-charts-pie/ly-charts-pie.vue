<template>
  <view class="ly-charts-pie" :style="{width: width}">
    <!-- 标题 -->
    <view 
      v-if="option.title && option.title.show !== false" 
      class="chart-title"
      :style="{
        textAlign: option.title.left === 'center' ? 'center' : option.title.left === 'right' ? 'right' : 'left'
      }"
    >
      <text 
        class="main-title"
        :style="{ color: option.title.textStyle?.color || '#333', fontSize: option.title.textStyle?.fontSize || 18 }"
      >{{ option.title.text }}</text>
      <text 
        v-if="option.title.subtext"
        class="sub-title"
        :style="{ color: option.title.subtextStyle?.color || '#666', fontSize: option.title.subtextStyle?.fontSize || 14 }"
      >{{ option.title.subtext }}</text>
    </view>
    
    <canvas 
      :id="cid" 
      :canvas-id="cid" 
      :style="{ width: canvasWidth + 'px', height: canvasHeight + 'px' }"
      @touchstart="tap"
      @touchmove="move"
      @touchend="touchEnd"
    ></canvas>
  </view>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick, getCurrentInstance } from 'vue';

const instance = getCurrentInstance()!.proxy!;

// #ifdef APP-NVUE
const dom = weex.requireModule('dom');
// #endif

// 定义组件属性
const props = defineProps({
  // 图表数据
  option: {
    type: Object,
    default: () => ({})
  },
  // 图表宽度
  width: {
    type: [String, Number],
    default: '100%'
  },
  // 图表高度
  height: {
    type: [String, Number],
    default: 300
  }
});

// 定义事件
const emit = defineEmits(['click', 'tooltipShow']);

// 生成唯一ID
const cid = 'ly-charts-pie-' + Math.random().toString(36).substr(2);

// 解析单位的辅助函数，支持rpx、px和数字
const parseUnit = (value) => {
  if (typeof value === 'number') {
    return value;
  }
  
  if (typeof value === 'string') {
    if (value.endsWith('rpx')) {
      return uni.upx2px(parseInt(value));
    } else if (value.endsWith('px')) {
      return parseInt(value);
    } else if (value.endsWith('%')) {
      return value; // 百分比保持原样
    } else {
      return parseInt(value) || 0;
    }
  }
  
  return 0;
};

// 响应式数据
const canvasWidth = ref(typeof props.width === 'string' && props.width.indexOf('%') !== -1 ? 
  null : 
  (typeof props.width === 'number' ? props.width : parseUnit(props.width)));
const canvasHeight = ref(typeof props.height === 'string' ? parseUnit(props.height) : props.height);
const isMount = ref(false);
const chartInstance = ref(null);

// 获取画布尺寸
const getCanvasSize = () => {
  return new Promise((resolve) => {
    // 如果canvasWidth已经为数字，直接返回
    if (typeof canvasWidth.value === 'number' && canvasWidth.value > 0) {
      resolve();
      return;
    }
    
    // 修复：在 setup 上下文中使用正确的 selector 查询方式
    const query = uni.createSelectorQuery().in(instance);
    query.select('#' + cid)
      .boundingClientRect((res) => {
        if (res) {
          canvasWidth.value = res.width || props.width;
          canvasHeight.value = res.height || props.height;
        } else {
          // 改进宽度计算逻辑，支持rpx单位
          canvasWidth.value = typeof props.width === 'string' && props.width.indexOf('%') !== -1 ? 
            (uni.upx2px(750) * parseInt(props.width) / 100) : 
            parseUnit(props.width);
          canvasHeight.value = parseUnit(props.height);
        }
        resolve();
      })
      .exec();
  });
};

// 绘制图表
const drawChart = () => {
  if (!props.option || !Object.keys(props.option).length) return;
  
  // 使用原生Canvas绘制饼图
  drawNativePie();
};

// 使用原生Canvas绘制饼图
const drawNativePie = () => {
  // #ifndef APP-NVUE
  // 修复：在 setup 上下文中正确创建 canvas context
  const ctx = uni.createCanvasContext(cid, instance);
  
  // 清空画布
  ctx.clearRect(0, 0, canvasWidth.value, canvasHeight.value);
  
  // 获取数据
  const series = props.option.series && props.option.series[0] || {};
  const data = series.data || [];
  if (!data.length) return;
  
  // 计算总值
  const total = data.reduce((sum, item) => sum + item.value, 0);
  
  // 设置饼图参数
  let centerX = canvasWidth.value / 2;
  let centerY = canvasHeight.value / 2;
  
  // 处理legend位置对饼图中心的影响
  let leftLegendWidth = 0;
  let rightLegendWidth = 0;
  let topLegendHeight = 0;
  let bottomLegendHeight = 0;
  
  if (props.option.legend && props.option.legend.show !== false) {
    // 简单估算图例尺寸
    const legendItemWidth = 80; // 假设每个图例项宽度为80px
    const legendItemHeight = 20; // 假设每个图例项高度为20px
    
    // 遵循ECharts规范，使用top、bottom、left、right控制图例位置
    // 支持数字（像素值）、百分比字符串或关键字
    let legend = props.option.legend;
    const topValue = legend.top;
    const bottomValue = legend.bottom;
    const leftValue = legend.left;
    const rightValue = legend.right;
    
    const hasTop = topValue !== undefined && topValue !== 'auto';
    const hasBottom = bottomValue !== undefined && bottomValue !== 'auto';
    const hasLeft = leftValue !== undefined && leftValue !== 'auto';
    const hasRight = rightValue !== undefined && rightValue !== 'auto';
    
    // 如果没有设置top或bottom，则默认放在底部
    const isBottomDefault = !hasTop && !hasBottom;
    
    if (props.option.legend.orient === 'vertical') {
      const legendHeight = data.length * legendItemHeight;
      if (hasLeft) {
        leftLegendWidth = legendItemWidth;
        centerX = (canvasWidth.value + leftLegendWidth) / 2;
      } else if (hasRight) {
        rightLegendWidth = legendItemWidth;
        centerX = (canvasWidth.value - rightLegendWidth) / 2;
      }
    } else { // horizontal
      const legendWidth = data.length * legendItemWidth;
      if (hasTop) {
        // 计算top值（支持数字和百分比）
        let topPos = 10;
        if (typeof topValue === 'number') {
          topPos = topValue;
        } else if (typeof topValue === 'string' && topValue.endsWith('%')) {
          topPos = (parseInt(topValue) / 100) * canvasHeight.value;
        }
        topLegendHeight = legendItemHeight * 2; // 给图例更多空间
        centerY = (canvasHeight.value + topLegendHeight) / 2;
      } else if (hasBottom || isBottomDefault) {
        // 计算bottom值（支持数字和百分比）
        let bottomPos = 10;
        if (typeof bottomValue === 'number') {
          bottomPos = bottomValue;
        } else if (typeof bottomValue === 'string' && bottomValue.endsWith('%')) {
          bottomPos = (parseInt(bottomValue) / 100) * canvasHeight.value;
        } else if (isBottomDefault) {
          bottomPos = 10; // 默认底部距离
        }
        bottomLegendHeight = legendItemHeight * 2;
        centerY = (canvasHeight.value - bottomLegendHeight) / 2;
      }
    }
  }
  
  // 处理半径，考虑图例占用空间
  let radiusValue = series.radius || '65%'; // 增大默认半径到70%
  
  // 支持内外半径数组格式 ['40%', '70%']
  let innerRadius, outerRadius;
  if (Array.isArray(radiusValue)) {
    innerRadius = radiusValue[0];
    outerRadius = radiusValue[1];
  } else {
    innerRadius = 0;
    outerRadius = radiusValue;
  }
  
  // 处理百分比格式
  if (typeof outerRadius === 'string' && outerRadius.endsWith('%')) {
    outerRadius = parseFloat(outerRadius) / 100;
  } else if (typeof outerRadius === 'number') {
    outerRadius = outerRadius / Math.min(canvasWidth.value, canvasHeight.value);
  }
  
  if (typeof innerRadius === 'string' && innerRadius.endsWith('%')) {
    innerRadius = parseFloat(innerRadius) / 100;
  } else if (typeof innerRadius === 'number') {
    innerRadius = innerRadius / Math.min(canvasWidth.value, canvasHeight.value);
  }
  
  // 计算可用于饼图的尺寸
  const availableWidth = canvasWidth.value - leftLegendWidth - rightLegendWidth - 20; // 减去边距
  const availableHeight = canvasHeight.value - topLegendHeight - bottomLegendHeight - 20;
  const outerRadiusPx = Math.min(availableWidth, availableHeight) * (typeof outerRadius === 'number' ? outerRadius : 0.7) / 2;
  const innerRadiusPx = outerRadiusPx * (innerRadius || 0);
  
  // 默认颜色列表
  const defaultColors = [
    '#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de',
    '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc'
  ];
  
  // 绘制饼图
  let startAngle = -Math.PI / 2; // 从顶部开始绘制
  const sectorAngles = []; // 保存每个扇形的角度信息
  
  for (let i = 0; i < data.length; i++) {
    const item = data[i];
    const percent = item.value / total;
    const endAngle = startAngle + percent * 2 * Math.PI;
    
    // 保存扇形角度信息
    sectorAngles.push({
      startAngle,
      endAngle,
      middleAngle: (startAngle + endAngle) / 2,
      percent,
      data: item
    });
    
    // 设置扇形样式
    ctx.beginPath();
    
    // 如果是环形图，绘制内外弧
    if (innerRadiusPx > 0) {
      // 绘制外弧
      ctx.arc(centerX, centerY, outerRadiusPx, startAngle, endAngle);
      // 绘制内弧（反向绘制）
      ctx.arc(centerX, centerY, innerRadiusPx, endAngle, startAngle, true);
    } else {
      // 绘制实心饼图
      ctx.moveTo(centerX, centerY);
      ctx.arc(centerX, centerY, outerRadiusPx, startAngle, endAngle);
      ctx.closePath();
    }
    
    ctx.setFillStyle(item.color || item.itemStyle?.color || defaultColors[i % defaultColors.length]);
    ctx.fill();
    
    // 绘制边框 (支持 itemStyle.borderWidth 和 itemStyle.borderColor)
    if (series.itemStyle?.borderWidth || item.itemStyle?.borderWidth) {
      const borderWidth = item.itemStyle?.borderWidth !== undefined ? 
        item.itemStyle.borderWidth : series.itemStyle.borderWidth;
      const borderColor = item.itemStyle?.borderColor || 
        series.itemStyle.borderColor || '#ffffff';
      
      ctx.setLineWidth(borderWidth);
      ctx.setStrokeStyle(borderColor);
      ctx.stroke();
    } else {
      // 默认边框
      ctx.setStrokeStyle('#ffffff');
      ctx.setLineWidth(2);
      ctx.stroke();
    }
    
    startAngle = endAngle;
  }
  
  // 绘制带引导线的标签
  drawLabelsWithLines(ctx, sectorAngles, centerX, centerY, outerRadiusPx, defaultColors, data);
  
  // 绘制legend
  if (props.option.legend && props.option.legend.show !== false) {
    drawLegend(ctx, data, defaultColors);
  }
  
  ctx.draw(false);
  
  // 保存图表实例信息
  chartInstance.value = {
    type: 'native-pie',
    data: data,
    centerX: centerX,
    centerY: centerY,
    radius: outerRadiusPx,
    innerRadius: innerRadiusPx,
    total: total,
    sectorAngles: sectorAngles, // 保存扇形角度信息用于点击判断
    destroy: () => {} // 空销毁函数
  };
  // #endif
};

// 绘制带引导线的标签
const drawLabelsWithLines = (ctx, sectorAngles, centerX, centerY, radius, defaultColors, data) => {
  // 标签和引导线相关参数
  const labelRadius = radius + 30; // 标签距离中心的距离
  const lineBreakRadius = radius + 15; // 引导线转折点距离中心的距离
  
  ctx.setFontSize(12);
  ctx.setTextBaseline('middle');
  
  for (let i = 0; i < sectorAngles.length; i++) {
    const sector = sectorAngles[i];
    const item = sector.data;
    const middleAngle = sector.middleAngle;
    
    // 计算扇形中点坐标（用于引导线起点）
    const startX = centerX + radius * Math.cos(middleAngle);
    const startY = centerY + radius * Math.sin(middleAngle);
    
    // 计算引导线转折点坐标
    const breakX = centerX + lineBreakRadius * Math.cos(middleAngle);
    const breakY = centerY + lineBreakRadius * Math.sin(middleAngle);
    
    // 计算标签位置（引导线终点）
    const endX = centerX + labelRadius * Math.cos(middleAngle);
    const endY = centerY + labelRadius * Math.sin(middleAngle);
    
    // 添加折弯效果 - 根据角度调整标签的x坐标
    let labelX = endX;
    let breakPointX = breakX;
    const horizontalOffset = 20; // 水平偏移量保持不变
    
    // 根据标签位置调整折弯方向
    if (Math.abs(middleAngle) < Math.PI / 2) {
      // 右侧标签 - 折弯向右
      labelX = endX + horizontalOffset;
      breakPointX = breakX + horizontalOffset;
    } else {
      // 左侧标签 - 折弯向左
      labelX = endX - horizontalOffset;
      breakPointX = breakX - horizontalOffset;
    }
    
    // 绘制引导线（带折弯）
    ctx.beginPath();
    ctx.moveTo(startX, startY);
    ctx.lineTo(breakX, breakY);
    ctx.lineTo(labelX, endY); // 现在末端是完全水平的
    ctx.setStrokeStyle('#666666');
    ctx.setLineWidth(1);
    ctx.stroke();
    
    // 绘制标签背景
    const text = item.name;
    ctx.setFontSize(12);
    const textWidth = ctx.measureText(text).width;
    
    // 根据标签位置调整对齐方式
    if (Math.abs(middleAngle) < Math.PI / 2) {
      // 右侧标签
      ctx.setFillStyle('#ffffff');
      ctx.fillRect(labelX, endY - 10, textWidth + 10, 20);
      ctx.setFillStyle('#333333');
      ctx.setTextAlign('left');
      ctx.fillText(text, labelX + 5, endY);
    } else {
      // 左侧标签
      ctx.setFillStyle('#ffffff');
      ctx.fillRect(labelX - textWidth - 10, endY - 10, textWidth + 10, 20);
      ctx.setFillStyle('#333333');
      ctx.setTextAlign('right');
      ctx.fillText(text, labelX - 5, endY);
    }
  }
};

// 绘制legend
const drawLegend = (ctx, data, defaultColors) => {
  const legend = props.option.legend;
  if (!legend || legend.show === false) return;
  
  const itemHeight = 14;
  const itemSpacing = 10;
  let startX, startY;
  
  // 遵循ECharts规范，使用top、bottom、left、right控制图例位置
  // 支持数字（像素值）、百分比字符串或关键字
  const topValue = legend.top;
  const bottomValue = legend.bottom;
  const leftValue = legend.left;
  const rightValue = legend.right;
  
  const hasTop = topValue !== undefined && topValue !== 'auto';
  const hasBottom = bottomValue !== undefined && bottomValue !== 'auto';
  const hasLeft = leftValue !== undefined && leftValue !== 'auto';
  const hasRight = rightValue !== undefined && rightValue !== 'auto';
  
  // 默认将图例放在底部
  const isBottomDefault = !hasTop && !hasBottom;
  
  if (legend.orient === 'vertical') {
    // 垂直排列
    const legendWidth = 80;
    const legendHeight = data.length * (itemHeight + itemSpacing);
    
    if (hasLeft) {
      // 计算left值（支持数字和百分比）
      if (typeof leftValue === 'number') {
        startX = leftValue;
      } else if (typeof leftValue === 'string' && leftValue.endsWith('%')) {
        startX = (parseInt(leftValue) / 100) * canvasWidth.value;
      } else {
        startX = 10;
      }
    } else if (hasRight) {
      // 计算right值（支持数字和百分比）
      let rightPos = 10;
      if (typeof rightValue === 'number') {
        rightPos = rightValue;
      } else if (typeof rightValue === 'string' && rightValue.endsWith('%')) {
        rightPos = (parseInt(rightValue) / 100) * canvasWidth.value;
      }
      startX = canvasWidth.value - legendWidth - rightPos;
    } else {
      startX = canvasWidth.value / 2 - legendWidth / 2;
    }
    
    if (hasTop) {
      // 计算top值（支持数字和百分比）
      if (typeof topValue === 'number') {
        startY = topValue;
      } else if (typeof topValue === 'string' && topValue.endsWith('%')) {
        startY = (parseInt(topValue) / 100) * canvasHeight.value;
      } else {
        startY = 10;
      }
    } else if (hasBottom || isBottomDefault) { // 默认在底部
      // 计算bottom值（支持数字和百分比）
      let bottomPos = 10;
      if (typeof bottomValue === 'number') {
        bottomPos = bottomValue;
      } else if (typeof bottomValue === 'string' && bottomValue.endsWith('%')) {
        bottomPos = (parseInt(bottomValue) / 100) * canvasHeight.value;
      } else if (isBottomDefault) {
        bottomPos = 10; // 默认底部距离
      }
      startY = canvasHeight.value - legendHeight - bottomPos;
    } else {
      startY = canvasHeight.value / 2 - legendHeight / 2;
    }
  } else {
    // 水平排列
    const legendWidth = data.length * 70; // 增加图例项宽度
    const legendHeight = itemHeight;
    
    if (hasLeft) {
      // 计算left值（支持数字和百分比）
      if (typeof leftValue === 'number') {
        startX = leftValue;
      } else if (typeof leftValue === 'string' && leftValue.endsWith('%')) {
        startX = (parseInt(leftValue) / 100) * canvasWidth.value;
      } else {
        startX = 10;
      }
    } else if (hasRight) {
      // 计算right值（支持数字和百分比）
      let rightPos = 10;
      if (typeof rightValue === 'number') {
        rightPos = rightValue;
      } else if (typeof rightValue === 'string' && rightValue.endsWith('%')) {
        rightPos = (parseInt(rightValue) / 100) * canvasWidth.value;
      }
      startX = canvasWidth.value - legendWidth - rightPos;
    } else {
      startX = canvasWidth.value / 2 - legendWidth / 2;
    }
    
    if (hasTop) {
      // 计算top值（支持数字和百分比）
      if (typeof topValue === 'number') {
        startY = topValue;
      } else if (typeof topValue === 'string' && topValue.endsWith('%')) {
        startY = (parseInt(topValue) / 100) * canvasHeight.value;
      } else {
        startY = 10;
      }
    } else if (hasBottom || isBottomDefault) { // 默认在底部
      // 计算bottom值（支持数字和百分比）
      let bottomPos = 10;
      if (typeof bottomValue === 'number') {
        bottomPos = bottomValue;
      } else if (typeof bottomValue === 'string' && bottomValue.endsWith('%')) {
        bottomPos = (parseInt(bottomValue) / 100) * canvasHeight.value;
      } else if (isBottomDefault) {
        bottomPos = 10; // 默认底部距离
      }
      startY = canvasHeight.value - legendHeight - bottomPos;
    } else {
      startY = 10; // 默认在顶部
    }
  }
  
  ctx.setFontSize(12);
  ctx.setTextBaseline('middle');
  
  for (let i = 0; i < data.length; i++) {
    const item = data[i];
    let x, y;
    
    if (legend.orient === 'vertical') {
      x = startX;
      y = startY + i * (itemHeight + itemSpacing);
    } else {
      x = startX + i * 70; // 增加图例项间距
      y = startY;
    }
    
    // 绘制颜色方块
    ctx.setFillStyle(item.color || defaultColors[i % defaultColors.length]);
    ctx.fillRect(x, y - 6, 12, 12);
    
    // 绘制文字
    ctx.setFillStyle('#333');
    ctx.setTextAlign('left');
    ctx.fillText(item.name, x + 18, y);
  }
};

// 点击事件
const tap = (e) => {
  if (chartInstance.value && chartInstance.value.type === 'native-pie') {
    // 原生饼图点击处理
    const { x, y } = e.touches[0];
    const dx = x - chartInstance.value.centerX;
    const dy = y - chartInstance.value.centerY;
    const distance = Math.sqrt(dx * dx + dy * dy);
    
    // 判断点击是否在饼图范围内
    if (distance <= chartInstance.value.radius) {
      const angle = (Math.atan2(dy, dx) + Math.PI / 2 + 2 * Math.PI) % (2 * Math.PI);
      let cumulativeAngle = 0;
      let clickedIndex = -1;
      
      // 查找点击的扇形
      for (let i = 0; i < chartInstance.value.data.length; i++) {
        const item = chartInstance.value.data[i];
        const sectorAngle = (item.value / chartInstance.value.total) * 2 * Math.PI;
        cumulativeAngle += sectorAngle;
        
        if (angle <= cumulativeAngle) {
          clickedIndex = i;
          break;
        }
      }
      
      if (clickedIndex !== -1) {
        const clickedData = chartInstance.value.data[clickedIndex];
        emit('click', {
          name: clickedData.name,
          value: clickedData.value,
          dataIndex: clickedIndex,
          data: clickedData,
          event: e
        });
        
        // 触发tooltip事件
        if (props.option.tooltip && props.option.tooltip.trigger === 'item') {
          emit('tooltipShow', {
            name: clickedData.name,
            value: clickedData.value,
            dataIndex: clickedIndex,
            data: clickedData
          });
        }
      }
    }
  }
};

// 移动事件
const move = (e) => {
  // 可以添加触摸移动逻辑
};

// 触摸结束事件
const touchEnd = (e) => {
  // 可以添加触摸结束逻辑
};

// 更新数据
const updateData = (data) => {
  if (chartInstance.value) {
    if (chartInstance.value.type === 'native-pie') {
      // 原生饼图更新数据
      props.option.series[0].data = data;
      drawNativePie();
    }
  }
};

// 获取图表实例
const getChartInstance = () => {
  return chartInstance.value;
};

// 监听option变化
watch(
  () => props.option,
  (newVal) => {
    if (isMount.value) {
      drawChart();
    }
  },
  { deep: true }
);

// 初始化
const init = () => {
  getCanvasSize().then(() => {
    drawChart();
  });
};

// 初始化nvue
const initNvue = () => {
  dom.getComponentRect(`#${cid}`, (res) => {
    if (res.result) {
      canvasWidth.value = res.size.width;
      canvasHeight.value = res.size.height;
      drawChart();
    }
  });
};

// 组件挂载
onMounted(() => {
  isMount.value = true;
  // #ifndef APP-NVUE
  nextTick(() => {
    init();
  });
  // #endif
  
  // #ifdef APP-NVUE
  initNvue();
  // #endif
});

// 组件卸载
onUnmounted(() => {
  if (chartInstance.value) {
    chartInstance.value.destroy();
  }
});

// 暴露方法给父组件
defineExpose({
  updateData,
  getChartInstance
});
</script>

<style scoped>
.ly-charts-pie {
  position: relative;
  width: 100%;
}

.chart-title {
  padding: 10px 0;
}

.main-title {
  display: block;
  font-weight: bold;
  margin-bottom: 5px;
}

.sub-title {
  display: block;
  opacity: 0.8;
}
</style>