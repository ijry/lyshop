<template>
  <view
    class="ly-canvas"
    :id="rootId"
    :style="{
      width: useRootHeightAndWidth ? '100%' : 'auto',
      height: useRootHeightAndWidth ? '100%' : 'auto'
    }"
  >
    <!-- #ifdef MP || H5 -->
    <canvas
      class="ly-canvas__canvas"
      :id="canvasId"
      :canvas-id="canvasId"
      type="2d"
      :style="{ width: actualWidth + unit, height: actualHeight + unit }"
      @touchstart="onTouchStart"
      @touchmove="onTouchMove"
      @touchend="onTouchEnd"
    />
    <!-- #endif -->

    <!-- #ifdef APP-PLUS -->
    <canvas
      class="ly-canvas__canvas"
      :id="canvasId"
      :canvas-id="canvasId"
      :style="{ width: actualWidth + unit, height: actualHeight + unit }"
      @touchstart="onTouchStart"
      @touchmove="onTouchMove"
      @touchend="onTouchEnd"
    />
    <!-- #endif -->

    <!-- #ifdef APP-NVUE -->
    <gcanvas
      class="ly-canvas__canvas"
      ref="gcanvas"
      :style="{ width: actualWidth + unit, height: actualHeight + unit }"
      @touchstart="onTouchStart"
      @touchmove="onTouchMove"
      @touchend="onTouchEnd"
    />
    <!-- #endif -->
  </view>
</template>

<script>
// #ifdef APP-NVUE
import { enable, WeexBridge } from '../../libs/util/gcanvas/index.js';
// #endif

let canvasNode = null;

export default {
  name: 'ly-canvas',
  emits: ['ready', 'touchstart', 'touchmove', 'touchend'],
  props: {
    canvasId: {
      type: String,
      default: () => `ly-canvas-${Math.floor(Math.random() * 1000000)}`
    },
    width: {
      type: [String, Number],
      default: 300
    },
    height: {
      type: [String, Number],
      default: 300
    },
    unit: {
      type: String,
      default: 'px'
    },
    useRootHeightAndWidth: {
      type: Boolean,
      default: false
    },
    bgColor: {
      type: String,
      default: 'transparent'
    }
  },
  data() {
    return {
      rootId: `ly-root-${Math.floor(Math.random() * 1000000)}`,
      ctx: null,
      widthLocal: this.parseSize(this.width),
      heightLocal: this.parseSize(this.height),
      fontSize: 12,
      fontFamily: 'sans-serif',
      fontWeight: 'normal'
    };
  },
  computed: {
    actualWidth() {
      return this.useRootHeightAndWidth ? this.widthLocal : this.parseSize(this.width);
    },
    actualHeight() {
      return this.useRootHeightAndWidth ? this.heightLocal : this.parseSize(this.height);
    }
  },
  methods: {
    parseSize(value) {
      if (typeof value === 'number') {
        return value;
      }
      if (typeof value !== 'string') {
        return 0;
      }
      if (value.endsWith('rpx') || value.endsWith('upx')) {
        return uni.upx2px(parseInt(value, 10));
      }
      if (value.endsWith('px')) {
        return parseInt(value, 10);
      }
      return parseInt(value, 10) || 0;
    },
    onTouchStart(event) {
      this.$emit('touchstart', event);
    },
    onTouchMove(event) {
      this.$emit('touchmove', event);
    },
    onTouchEnd(event) {
      this.$emit('touchend', event);
    },
    getNode(id, isCanvas = true) {
      return new Promise((resolve) => {
        try {
          // #ifdef APP-NVUE
          setTimeout(() => {
            const gcanvas = this.$refs.gcanvas;
            if (!gcanvas) {
              resolve(false);
              return;
            }
            resolve(enable(gcanvas, { bridge: WeexBridge }));
          }, 100);
          // #endif

          // #ifndef APP-PLUS-NVUE
          uni.createSelectorQuery()
            .in(this)
            .select(`#${id}`)
            .fields(
              {
                node: isCanvas,
                size: true
              },
              (res) => {
                resolve(res || false);
              }
            )
            .exec();
          // #endif
        } catch (error) {
          console.error('获取画布节点失败:', error);
          resolve(false);
        }
      });
    },
    getCanvasContext() {
      // #ifdef APP-PLUS
      return uni.createCanvasContext(this.canvasId, this);
      // #endif
      // #ifdef APP-PLUS-NVUE || MP || H5
      return canvasNode ? canvasNode.getContext('2d') : null;
      // #endif
    },
    applyFont() {
      if (!this.ctx) {
        return;
      }
      const font = `${this.fontWeight === 'normal' ? '' : `${this.fontWeight} `}${this.fontSize}px ${this.fontFamily}`.trim();
      if (typeof this.ctx.setFont === 'function') {
        this.ctx.setFont(font);
      } else if ('font' in this.ctx) {
        this.ctx.font = font;
      } else if (typeof this.ctx.setFontSize === 'function') {
        this.ctx.setFontSize(this.fontSize);
      }
    },
    async setNewSize() {
      const rootNode = await this.getNode(this.rootId, false);
      if (!rootNode) {
        return;
      }
      if (rootNode.width) {
        this.widthLocal = rootNode.width;
      }
      if (rootNode.height) {
        this.heightLocal = rootNode.height;
      }
    },
    async initCanvas(force = false) {
      if (this.useRootHeightAndWidth) {
        await this.setNewSize();
      }
      if (this.ctx && !force) {
        this.$emit('ready', {
          width: this.actualWidth,
          height: this.actualHeight
        });
        return true;
      }

      canvasNode = await this.getNode(this.canvasId);
      if (!canvasNode) {
        return false;
      }

      // #ifdef MP-WEIXIN
      const dpr = uni.getSystemInfoSync().pixelRatio || 1;
      canvasNode.width = this.actualWidth * dpr;
      canvasNode.height = this.actualHeight * dpr;
      // #endif

      this.ctx = this.getCanvasContext();
      if (!this.ctx) {
        return false;
      }

      // #ifdef MP-WEIXIN
      this.ctx.scale(dpr, dpr);
      // #endif

      this.applyFont();
      this.clearCanvas();
      this.$emit('ready', {
        width: this.actualWidth,
        height: this.actualHeight
      });
      return true;
    },
    refresh() {
      return this.initCanvas(true);
    },
    getWidth() {
      return this.actualWidth;
    },
    getHeight() {
      return this.actualHeight;
    },
    getRawContext() {
      return this.ctx;
    },
    clearCanvas() {
      if (!this.ctx) {
        return;
      }
      this.clearRect(0, 0, this.actualWidth, this.actualHeight);
      if (this.bgColor && this.bgColor !== 'transparent') {
        this.beginPath();
        this.rect(0, 0, this.actualWidth, this.actualHeight);
        this.setFillStyle(this.bgColor);
        this.fill();
      }
      this.draw();
    },
    rect(x, y, width, height) {
      if (this.ctx) {
        this.ctx.rect(x, y, width, height);
      }
    },
    clearRect(x, y, width, height) {
      if (this.ctx) {
        this.ctx.clearRect(x, y, width, height);
      }
    },
    fillRect(x, y, width, height) {
      if (this.ctx) {
        this.ctx.fillRect(x, y, width, height);
      }
    },
    strokeRect(x, y, width, height) {
      if (this.ctx) {
        this.ctx.strokeRect(x, y, width, height);
      }
    },
    fill() {
      if (this.ctx) {
        this.ctx.fill();
      }
    },
    stroke() {
      if (this.ctx) {
        this.ctx.stroke();
      }
    },
    beginPath() {
      if (this.ctx) {
        this.ctx.beginPath();
      }
    },
    closePath() {
      if (this.ctx) {
        this.ctx.closePath();
      }
    },
    moveTo(x, y) {
      if (this.ctx) {
        this.ctx.moveTo(x, y);
      }
    },
    lineTo(x, y) {
      if (this.ctx) {
        this.ctx.lineTo(x, y);
      }
    },
    arc(x, y, radius, startAngle, endAngle, anticlockwise = false) {
      if (this.ctx) {
        this.ctx.arc(x, y, radius, startAngle, endAngle, anticlockwise);
      }
    },
    bezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y) {
      if (this.ctx) {
        this.ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y);
      }
    },
    quadraticCurveTo(cpx, cpy, x, y) {
      if (this.ctx) {
        this.ctx.quadraticCurveTo(cpx, cpy, x, y);
      }
    },
    save() {
      if (this.ctx && typeof this.ctx.save === 'function') {
        this.ctx.save();
      }
    },
    restore() {
      if (this.ctx && typeof this.ctx.restore === 'function') {
        this.ctx.restore();
      }
    },
    translate(x, y) {
      if (this.ctx && typeof this.ctx.translate === 'function') {
        this.ctx.translate(x, y);
      }
    },
    rotate(angle) {
      if (this.ctx && typeof this.ctx.rotate === 'function') {
        this.ctx.rotate(angle);
      }
    },
    scale(x, y) {
      if (this.ctx && typeof this.ctx.scale === 'function') {
        this.ctx.scale(x, y);
      }
    },
    setFillStyle(color) {
      if (!this.ctx) {
        return;
      }
      if (typeof this.ctx.setFillStyle === 'function') {
        this.ctx.setFillStyle(color);
      } else {
        this.ctx.fillStyle = color;
      }
    },
    setStrokeStyle(color) {
      if (!this.ctx) {
        return;
      }
      if (typeof this.ctx.setStrokeStyle === 'function') {
        this.ctx.setStrokeStyle(color);
      } else {
        this.ctx.strokeStyle = color;
      }
    },
    setLineWidth(width) {
      if (!this.ctx) {
        return;
      }
      if (typeof this.ctx.setLineWidth === 'function') {
        this.ctx.setLineWidth(width);
      } else {
        this.ctx.lineWidth = width;
      }
    },
    setLineCap(lineCap = 'round') {
      if (!this.ctx) {
        return;
      }
      if (typeof this.ctx.setLineCap === 'function') {
        this.ctx.setLineCap(lineCap);
      } else {
        this.ctx.lineCap = lineCap;
      }
    },
    setLineJoin(lineJoin = 'round') {
      if (!this.ctx) {
        return;
      }
      if (typeof this.ctx.setLineJoin === 'function') {
        this.ctx.setLineJoin(lineJoin);
      } else {
        this.ctx.lineJoin = lineJoin;
      }
    },
    setTextAlign(align = 'left') {
      if (!this.ctx) {
        return;
      }
      if (typeof this.ctx.setTextAlign === 'function') {
        this.ctx.setTextAlign(align);
      } else {
        this.ctx.textAlign = align;
      }
    },
    setTextBaseline(baseline = 'alphabetic') {
      if (!this.ctx) {
        return;
      }
      if (typeof this.ctx.setTextBaseline === 'function') {
        this.ctx.setTextBaseline(baseline);
      } else {
        this.ctx.textBaseline = baseline;
      }
    },
    setFontSize(fontSize) {
      this.fontSize = fontSize;
      this.applyFont();
    },
    setFont(font) {
      if (!this.ctx) {
        return;
      }
      if ('font' in this.ctx) {
        this.ctx.font = font;
        return;
      }
      const matched = String(font).match(/(\d+)px/);
      if (matched) {
        this.fontSize = Number(matched[1]);
      }
      this.applyFont();
    },
    setShadow(offsetX = 0, offsetY = 0, blur = 0, color = 'rgba(0,0,0,0)') {
      if (!this.ctx) {
        return;
      }
      this.ctx.shadowOffsetX = offsetX;
      this.ctx.shadowOffsetY = offsetY;
      this.ctx.shadowBlur = blur;
      this.ctx.shadowColor = color;
    },
    fillText(text, x, y) {
      if (this.ctx) {
        this.ctx.fillText(String(text), x, y);
      }
    },
    measureText(text) {
      if (this.ctx && typeof this.ctx.measureText === 'function') {
        return this.ctx.measureText(String(text));
      }
      return {
        width: String(text).length * this.fontSize * 0.6
      };
    },
    createLinearGradient(x0, y0, x1, y1) {
      if (this.ctx && typeof this.ctx.createLinearGradient === 'function') {
        return this.ctx.createLinearGradient(x0, y0, x1, y1);
      }
      return null;
    },
    draw() {
      if (this.ctx && typeof this.ctx.draw === 'function') {
        this.ctx.draw();
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

<style lang="scss" scoped>
.ly-canvas {
  position: relative;
  overflow: hidden;
}

.ly-canvas__canvas {
  display: block;
  width: 100%;
  height: 100%;
}
</style>
