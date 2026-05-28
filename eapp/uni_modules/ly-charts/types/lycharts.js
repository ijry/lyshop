// UTS兼容的图表选项类型定义
type LYCHartsOption = {
  // 背景色
  backgroundColor?: string;
  
  // 标题配置
  title?: LYCHartsTitleConfig;
  
  // 图例配置
  legend?: LYCHartsLegendConfig;
  
  // X轴配置
  xAxis?: LYCHartsAxisConfig;
  
  // Y轴配置
  yAxis?: LYCHartsAxisConfig;
  
  // 系列数据配置
  series: LYCHartsSeriesItem[];
  
  // 网格配置
  grid?: LYCHartsGridConfig | boolean;
  
  // 动画配置
  animation?: boolean;
  animationDuration?: number;
  
  // 提示框配置
  tooltip?: LYCHartsTooltipConfig;
  
  // 工具栏配置
  toolbox?: LYCHartsToolboxConfig;
  
  // 数据缩放配置
  dataZoom?: LYCHartsDataZoomConfig[];
}

// 标题配置
type LYCHartsTitleConfig = {
  show?: boolean;
  text?: string;
  subtext?: string;
  left?: 'left' | 'center' | 'right' | string;
  top?: 'top' | 'middle' | 'bottom' | string;
  textStyle?: LYCHartsTextStyle;
  subtextStyle?: LYCHartsTextStyle;
};

// 文字样式
type LYCHartsTextStyle = {
  color?: string;
  fontSize?: number;
  fontWeight?: string | number;
};

// 图例配置
type LYCHartsLegendConfig = {
  show?: boolean;
  data: string[];
  left?: 'left' | 'center' | 'right' | string;
  top?: 'top' | 'middle' | 'bottom' | string;
  orient?: 'horizontal' | 'vertical';
};

// 坐标轴配置
type LYCHartsAxisConfig = {
  show?: boolean;
  type?: 'value' | 'category' | 'time' | 'log';
  name?: string;
  nameLocation?: 'start' | 'middle' | 'end';
  nameGap?: number;
  nameTextStyle?: LYCHartsTextStyle;
  min?: number | 'dataMin' | ((value: { min: number }) => number);
  max?: number | 'dataMax' | ((value: { max: number }) => number);
  interval?: number;
  axisLine?: LYCHartsAxisLineStyle;
  axisTick?: LYCHartsAxisTickStyle;
  axisLabel?: LYCHartsAxisLabelStyle;
  splitLine?: LYCHartsSplitLineStyle;
  data?: any[];
};

// 轴线样式
type LYCHartsAxisLineStyle = {
  show?: boolean;
  lineStyle?: LYCHartsLineStyle;
};

// 轴刻度样式
type LYCHartsAxisTickStyle = {
  show?: boolean;
  alignWithLabel?: boolean;
};

// 轴标签样式
type LYCHartsAxisLabelStyle = {
  show?: boolean;
  rotate?: number;
  margin?: number;
};

// 分割线样式
type LYCHartsSplitLineStyle = {
  show?: boolean;
  lineStyle?: LYCHartsLineStyle;
};

// 系列项配置
type LYCHartsSeriesItem = {
  name?: string;
  type: 'line' | 'bar' | 'pie' | 'scatter' | 'gauge' | 'radar';
  data: any[];
  smooth?: boolean;
  stack?: string;
  xAxisIndex?: number;
  yAxisIndex?: number;
  color?: string;
  barWidth?: number | string;
  barGap?: string;
  itemStyle?: LYCHartsItemStyle;
  label?: LYCHartsLabelConfig;
  emphasis?: LYCHartsEmphasisConfig;
  radius?: string | number;  // 饼图用
  center?: [string | number, string | number];  // 饼图用
  startAngle?: number;  // 饼图用
  endAngle?: number;  // 饼图用
  min?: number;  // 仪表盘用
  max?: number;  // 仪表盘用
  splitNumber?: number;  // 仪表盘用
  pointer?: LYCHartsPointerConfig;  // 仪表盘用
  axisLine?: LYCHartsAxisLineStyle;  // 仪表盘用
  detail?: LYCHartsDetailConfig;  // 仪表盘用
};

// 项目样式
type LYCHartsItemStyle = {
  color?: string;
  borderColor?: string;
  borderWidth?: number;
  borderType?: 'solid' | 'dashed' | 'dotted';
  shadowBlur?: number;
  shadowColor?: string;
  shadowOffsetX?: number;
  shadowOffsetY?: number;
};

// 标签配置
type LYCHartsLabelConfig = {
  show?: boolean;
  position?: 'top' | 'inside' | 'outside' | 'center' | 'left' | 'right' | 'bottom';
  formatter?: string | ((params: any) => string);
  color?: string;
  fontSize?: number;
};

// 强调样式
type LYCHartsEmphasisConfig = {
  itemStyle?: LYCHartsItemStyle;
  label?: LYCHartsLabelConfig;
};

// 线条样式
type LYCHartsLineStyle = {
  color?: string;
  width?: number;
  type?: 'solid' | 'dashed' | 'dotted';
};

// 指针配置（主要用于仪表盘）
type LYCHartsPointerConfig = {
  show?: boolean;
  length?: string;
  width?: number;
};

// 详情配置（主要用于仪表盘）
type LYCHartsDetailConfig = {
  show?: boolean;
  offsetCenter?: [string, string];
  color?: string;
  fontSize?: number;
  formatter?: string | ((value: number) => string);
};

// 网格配置
type LYCHartsGridConfig = {
  show?: boolean;
  left?: number | string;
  right?: number | string;
  top?: number | string;
  bottom?: number | string;
  backgroundColor?: string;
  borderColor?: string;
  borderWidth?: number;
  containLabel?: boolean;
};

// 提示框配置
type LYCHartsTooltipConfig = {
  show?: boolean;
  trigger?: 'item' | 'axis' | 'none';
  formatter?: string | ((params: any) => string);
  backgroundColor?: string;
  borderColor?: string;
  borderWidth?: number;
  textStyle?: LYCHartsTextStyle;
};

// 工具栏配置
type LYCHartsToolboxConfig = {
  show?: boolean;
  feature?: LYCHartsToolboxFeature;
};

// 工具栏功能
type LYCHartsToolboxFeature = {
  saveAsImage?: LYCHartsToolboxSaveImage;
  dataView?: LYCHartsToolboxDataView;
  dataZoom?: LYCHartsToolboxDataZoom;
  magicType?: LYCHartsToolboxMagicType;
};

// 工具栏保存图片功能
type LYCHartsToolboxSaveImage = {
  show?: boolean;
  title?: string;
  type?: 'png' | 'jpeg';
  name?: string;
};

// 工具栏数据视图功能
type LYCHartsToolboxDataView = {
  show?: boolean;
  title?: string;
  lang?: string[];
};

// 工具栏数据缩放功能
type LYCHartsToolboxDataZoom = {
  show?: boolean;
  title?: { zoom?: string; back?: string };
};

// 工具栏动态类型功能
type LYCHartsToolboxMagicType = {
  show?: boolean;
  type?: string[];
  title?: { line?: string; bar?: string; pie?: string };
};

// 数据缩放配置
type LYCHartsDataZoomConfig = {
  type?: 'slider' | 'inside';
  show?: boolean;
  start?: number;
  end?: number;
  xAxisIndex?: number | number[];
  yAxisIndex?: number | number[];
};