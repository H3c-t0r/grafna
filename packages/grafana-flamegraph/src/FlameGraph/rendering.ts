import uFuzzy from '@leeoniya/ufuzzy';
import { RefObject, useEffect, useMemo, useState } from 'react';
import color from 'tinycolor2';

import { GrafanaTheme2 } from '@grafana/data';

import {
  BAR_BORDER_WIDTH,
  BAR_TEXT_PADDING_LEFT,
  COLLAPSE_THRESHOLD,
  HIDE_THRESHOLD,
  LABEL_THRESHOLD,
  PIXELS_PER_LEVEL,
} from '../constants';
import { ClickedItemData, ColorScheme, ColorSchemeDiff, TextAlign } from '../types';

import { getBarColorByDiff, getBarColorByPackage, getBarColorByValue } from './colors';
import { FlameGraphDataContainer, LevelItem } from './dataTransform';

const ufuzzy = new uFuzzy();

type RenderOptions = {
  canvasRef: RefObject<HTMLCanvasElement>;
  data: FlameGraphDataContainer;
  root: LevelItem;
  direction: 'children' | 'parents';

  // Depth in number of levels
  depth: number;
  wrapperWidth: number;

  // If we are rendering only zoomed in part of the graph.
  rangeMin: number;
  rangeMax: number;

  search: string;
  textAlign: TextAlign;

  // Total ticks that will be used for sizing
  totalViewTicks: number;
  // Total ticks that will be used for computing colors as some color scheme (like in diff view) should not be affected
  // by sandwich or focus view.
  totalColorTicks: number;
  // Total ticks used to compute the diff colors
  totalTicksRight: number | undefined;
  colorScheme: ColorScheme | ColorSchemeDiff;
  focusedItemData?: ClickedItemData;
  getTheme: () => GrafanaTheme2;
};

export function useFlameRender(options: RenderOptions) {
  const {
    canvasRef,
    data,
    root,
    depth,
    direction,
    wrapperWidth,
    rangeMin,
    rangeMax,
    search,
    textAlign,
    totalViewTicks,
    totalColorTicks,
    totalTicksRight,
    colorScheme,
    focusedItemData,
    getTheme,
  } = options;
  const foundLabels = useMemo(() => {
    if (search) {
      const foundLabels = new Set<string>();
      let idxs = ufuzzy.filter(data.getUniqueLabels(), search);

      if (idxs) {
        for (let idx of idxs) {
          foundLabels.add(data.getUniqueLabels()[idx]);
        }
      }

      return foundLabels;
    }
    // In this case undefined means there was no search so no attempt to highlighting anything should be made.
    return undefined;
  }, [search, data]);

  const ctx = useSetupCanvas(canvasRef, wrapperWidth, depth);
  const theme = getTheme();

  useEffect(() => {
    if (!ctx) {
      return;
    }

    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
    const pixelsPerTick = (wrapperWidth * window.devicePixelRatio) / totalViewTicks / (rangeMax - rangeMin);

    const stack: Array<{ item: LevelItem }> = [];
    stack.push({ item: root });

    while (stack.length > 0) {
      const args = stack.shift()!;
      const item = args.item;

      const barX = getBarX(item.start, totalViewTicks, rangeMin, pixelsPerTick);
      const barY = item.level * PIXELS_PER_LEVEL;

      let curBarTicks = item.value;
      const collapsed = curBarTicks * pixelsPerTick <= COLLAPSE_THRESHOLD;
      const width = curBarTicks * pixelsPerTick - (collapsed ? 0 : BAR_BORDER_WIDTH * 2);

      if (width < HIDE_THRESHOLD) {
        // We don't render nor it's children
        continue;
      }

      const height = PIXELS_PER_LEVEL;

      ctx.beginPath();
      ctx.rect(barX + (collapsed ? 0 : BAR_BORDER_WIDTH), barY, width, height);

      let label = data.getLabel(item.itemIndexes[0]);

      ctx.fillStyle = getFillStyle(
        colorScheme,
        item,
        totalColorTicks,
        totalTicksRight,
        theme,
        label,
        rangeMin,
        rangeMax,
        foundLabels,
        collapsed,
        focusedItemData ? focusedItemData.item.level : 0
      );

      if (collapsed) {
        // Only fill the collapsed rects
        ctx.fill();
      } else {
        ctx.stroke();
        ctx.fill();

        if (width >= LABEL_THRESHOLD) {
          renderLabel(ctx, data, label, item, width, barX, barY, textAlign);
        }
      }

      const nextList = direction === 'children' ? item.children : item.parents;
      if (nextList) {
        stack.unshift(...nextList.map((c) => ({ item: c })));
      }
    }
  }, [
    ctx,
    data,
    root,
    depth,
    wrapperWidth,
    rangeMin,
    rangeMax,
    search,
    focusedItemData,
    foundLabels,
    textAlign,
    totalViewTicks,
    totalColorTicks,
    totalTicksRight,
    colorScheme,
    theme,
    direction,
  ]);
}

function getFillStyle(
  colorScheme: ColorScheme | ColorSchemeDiff,
  item: LevelItem,
  totalTicks: number,
  totalTicksRight: number | undefined,
  theme: GrafanaTheme2,
  label: string,
  rangeMin: number,
  rangeMax: number,
  foundNames: Set<string> | undefined,
  collapsed: boolean,
  topLevel: number
) {
  const barColor =
    item.valueRight !== undefined &&
    (colorScheme === ColorSchemeDiff.Default || colorScheme === ColorSchemeDiff.DiffColorBlind)
      ? getBarColorByDiff(item.value, item.valueRight!, totalTicks, totalTicksRight!, colorScheme)
      : colorScheme === ColorScheme.ValueBased
      ? getBarColorByValue(item.value, totalTicks, rangeMin, rangeMax)
      : getBarColorByPackage(label, theme);

  const barMutedColor = color(theme.colors.background.secondary);
  const barMutedColorHex = theme.isLight
    ? barMutedColor.darken(10).toHexString()
    : barMutedColor.lighten(10).toHexString();

  if (foundNames) {
    // Means we are searching, we use color for matches and gray the rest
    return foundNames.has(label) ? barColor.toHslString() : barMutedColorHex;
  }

  // No search
  if (collapsed) {
    // Collapsed are always grayed
    return barMutedColorHex;
  } else {
    // Mute if we are above the focused symbol
    return item.level > topLevel - 1 ? barColor.toHslString() : barColor.lighten(15).toHslString();
  }
}

function useSetupCanvas(canvasRef: RefObject<HTMLCanvasElement>, wrapperWidth: number, numberOfLevels: number) {
  const [ctx, setCtx] = useState<CanvasRenderingContext2D>();

  useEffect(() => {
    if (!(numberOfLevels && canvasRef.current)) {
      return;
    }
    const ctx = canvasRef.current.getContext('2d')!;

    const height = PIXELS_PER_LEVEL * numberOfLevels;
    canvasRef.current.width = Math.round(wrapperWidth * window.devicePixelRatio);
    canvasRef.current.height = Math.round(height);
    canvasRef.current.style.width = `${wrapperWidth}px`;
    canvasRef.current.style.height = `${height / window.devicePixelRatio}px`;

    ctx.textBaseline = 'middle';
    ctx.font = 12 * window.devicePixelRatio + 'px monospace';
    ctx.strokeStyle = 'white';
    setCtx(ctx);
  }, [canvasRef, setCtx, wrapperWidth, numberOfLevels]);
  return ctx;
}

// Renders a text inside the node rectangle. It allows setting alignment of the text left or right which takes effect
// when text is too long to fit in the rectangle.
function renderLabel(
  ctx: CanvasRenderingContext2D,
  data: FlameGraphDataContainer,
  label: string,
  item: LevelItem,
  width: number,
  x: number,
  y: number,
  textAlign: TextAlign
) {
  ctx.save();
  ctx.clip(); // so text does not overflow
  ctx.fillStyle = '#222';

  const displayValue = data.valueDisplayProcessor(item.value);
  const unit = displayValue.suffix ? displayValue.text + displayValue.suffix : displayValue.text;

  // We only measure name here instead of full label because of how we deal with the units and aligning later.
  const measure = ctx.measureText(label);
  const spaceForTextInRect = width - BAR_TEXT_PADDING_LEFT;

  let fullLabel = `${label} (${unit})`;
  let labelX = Math.max(x, 0) + BAR_TEXT_PADDING_LEFT;

  // We use the desired alignment only if there is not enough space for the text, otherwise we keep left alignment as
  // that will already show full text.
  if (measure.width > spaceForTextInRect) {
    ctx.textAlign = textAlign;
    // If aligned to the right we don't want to take the space with the unit label as the assumption is user wants to
    // mainly see the name. This also reflects how pyro/flamegraph works.
    if (textAlign === 'right') {
      fullLabel = label;
      labelX = x + width - BAR_TEXT_PADDING_LEFT;
    }
  }

  ctx.fillText(fullLabel, labelX, y + PIXELS_PER_LEVEL / 2);
  ctx.restore();
}

/**
 * Returns the X position of the bar. totalTicks * rangeMin is to adjust for any current zoom. So if we zoom to a
 * section of the graph we align and shift the X coordinates accordingly.
 * @param offset
 * @param totalTicks
 * @param rangeMin
 * @param pixelsPerTick
 */
export function getBarX(offset: number, totalTicks: number, rangeMin: number, pixelsPerTick: number) {
  return (offset - totalTicks * rangeMin) * pixelsPerTick;
}
