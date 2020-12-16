import tinycolor from 'tinycolor2';
import uPlot, { Series } from 'uplot';
import { getCanvasContext } from '../../../utils/measureText';
import {
  DrawStyle,
  LineConfig,
  AreaConfig,
  PointsConfig,
  PointVisibility,
  LineInterpolation,
  AreaGradientMode,
} from '../config';
import { PlotConfigBuilder } from '../types';

export interface SeriesProps extends LineConfig, AreaConfig, PointsConfig {
  drawStyle: DrawStyle;
  scaleKey: string;
  height: number;
}

export class UPlotSeriesBuilder extends PlotConfigBuilder<SeriesProps, Series> {
  getConfig() {
    const {
      drawStyle,
      lineInterpolation,
      lineColor,
      lineWidth,
      showPoints,
      pointColor,
      pointSize,
      fillColor,
      fillOpacity,
      fillGradient,
      scaleKey,
      spanNulls,
    } = this.props;

    let lineConfig: Partial<Series> = {};

    if (drawStyle === DrawStyle.Points) {
      lineConfig.paths = () => null;
    } else {
      lineConfig.stroke = lineColor;
      lineConfig.width = lineWidth;
      lineConfig.paths = (self: uPlot, seriesIdx: number, idx0: number, idx1: number) => {
        let pathsBuilder = mapDrawStyleToPathBuilder(drawStyle, lineInterpolation);
        return pathsBuilder(self, seriesIdx, idx0, idx1);
      };
    }

    const pointsConfig: Partial<Series> = {
      points: {
        stroke: pointColor,
        fill: pointColor,
        size: pointSize,
      },
    };

    // we cannot set points.show property above (even to undefined) as that will clear uPlot's default auto behavior
    if (showPoints === PointVisibility.Auto) {
      if (drawStyle === DrawStyle.Bars) {
        pointsConfig.points!.show = false;
      }
    } else if (showPoints === PointVisibility.Never) {
      pointsConfig.points!.show = false;
    } else if (showPoints === PointVisibility.Always) {
      pointsConfig.points!.show = true;
    }

    let fillConfig: any | undefined;
    let fillOpacityNumber = fillOpacity ?? 0;

    if (fillColor) {
      fillConfig = {
        fill: fillColor,
      };
    }

    if (fillOpacityNumber !== 0) {
      fillConfig = {
        fill: tinycolor(fillColor ?? lineColor)
          .setAlpha(fillOpacityNumber / 100)
          .toRgbString(),
      };

      if (fillGradient && fillGradient !== AreaGradientMode.None) {
        fillConfig.fill = getCanvasGradient(fillConfig.fill, fillGradient, fillOpacityNumber);
      }
    }

    return {
      scale: scaleKey,
      spanGaps: spanNulls,
      ...lineConfig,
      ...pointsConfig,
      ...fillConfig,
    };
  }
}

interface PathBuilders {
  bars: Series.PathBuilder;
  linear: Series.PathBuilder;
  smooth: Series.PathBuilder;
  stepBefore: Series.PathBuilder;
  stepAfter: Series.PathBuilder;
}

let builders: PathBuilders | undefined = undefined;

function mapDrawStyleToPathBuilder(
  style: DrawStyle,
  lineInterpolation?: LineInterpolation,
  opts?: any
): Series.PathBuilder {
  // This should be global static, but Jest initalization was failing so we lazy load to avoid the issue
  if (!builders) {
    const pathBuilders = uPlot.paths;
    const barWidthFactor = 0.6;
    const barMaxWidth = Infinity;

    builders = {
      bars: pathBuilders.bars!({ size: [barWidthFactor, barMaxWidth] }),
      linear: pathBuilders.linear!(),
      smooth: pathBuilders.spline!(),
      stepBefore: pathBuilders.stepped!({ align: -1 }),
      stepAfter: pathBuilders.stepped!({ align: 1 }),
    };
  }

  if (style === DrawStyle.Bars) {
    return builders.bars;
  }
  if (style === DrawStyle.Line) {
    if (lineInterpolation === LineInterpolation.StepBefore) {
      return builders.stepBefore;
    }
    if (lineInterpolation === LineInterpolation.StepAfter) {
      return builders.stepAfter;
    }
    if (lineInterpolation === LineInterpolation.Smooth) {
      return builders.smooth;
    }
  }

  return builders.linear; // the default
}

function getCanvasGradient(color: string, gradientMode: AreaGradientMode, opacity: number): CanvasGradient {
  const ctx = getCanvasContext();
  const pxRatio = window.devicePixelRatio ?? 1;
  const pixelHeight = (plotHeight - 30) * pxRatio;
  const gradient = ctx.createLinearGradient(0, 0, 0, pixelHeight);

  switch (gradientMode) {
    case AreaGradientMode.Hue:
      const color1 = tinycolor(color)
        .spin(-25)
        .darken(30)
        .setAlpha(opacity)
        .toRgbString();
      const color2 = tinycolor(color)
        .spin(25)
        .lighten(35)
        .setAlpha(opacity)
        .toRgbString();
      gradient.addColorStop(0, color2);
      gradient.addColorStop(1, color1);

    case AreaGradientMode.Opacity:
    default:
      gradient.addColorStop(
        0,
        tinycolor(color)
          .setAlpha(opacity)
          .toRgbString()
      );
      gradient.addColorStop(
        1,
        tinycolor(color)
          .setAlpha(0)
          .toRgbString()
      );
      return gradient;
  }
}
