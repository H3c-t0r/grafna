import uPlot, { Series, Cursor } from 'uplot';
import { FIXED_UNIT } from '../GraphNG/GraphNG';
import { Quadtree, Rect, pointWithin } from '../BarChart/quadtree';
import { distribute, SPACE_BETWEEN } from '../BarChart/distribute';
import { TimelineMode } from './types';

const { round, min } = Math;

const pxRatio = devicePixelRatio;

const laneDistr = SPACE_BETWEEN;

const font = Math.round(10 * pxRatio) + 'px Arial';

type WalkCb = (idx: number, offPx: number, dimPx: number) => void;

function walk(laneWidth: number, yIdx: number | null, count: number, dim: number, draw: WalkCb) {
  distribute(count, laneWidth, laneDistr, yIdx, (i, offPct, dimPct) => {
    let laneOffPx = dim * offPct;
    let laneWidPx = dim * dimPct;

    draw(i, laneOffPx, laneWidPx);
  });
}

/**
 * @internal
 */
export interface TimelineCoreOptions {
  mode: TimelineMode;
  count: number;
  laneWidth: number;

  /** used only for grid mode */
  align?: -1 | 0 | 1;
  size?: [number, number?]; //[factor: number, max?: number];

  label: (seriesIdx: number) => string;
  fill: (seriesIdx: number, valueIdx: number, value: any) => CanvasRenderingContext2D['fillStyle'];
  stroke: (seriesIdx: number, valueIdx: number, value: any) => CanvasRenderingContext2D['strokeStyle'];
  formatValue?: (seriesIdx: number, value: any) => string;
  onHover?: (seriesIdx: number, valueIdx: number) => void;
  onLeave?: (seriesIdx: number, valueIdx: number) => void;
}

/**
 * @internal
 */
export function getConfig(opts: TimelineCoreOptions) {
  const {
    mode,
    count,
    laneWidth,
    align = 0,
    size = [0.6, Infinity],
    label,
    fill,
    stroke,
    formatValue,
    // onHover,
    // onLeave,
  } = opts;

  let qt: Quadtree;

  const hoverMarks = Array(count)
    .fill(null)
    .map(() => {
      let mark = document.createElement('div');
      mark.classList.add('bar-mark');
      mark.style.position = 'absolute';
      mark.style.background = 'rgba(255,255,255,0.4)';
      return mark;
    });

  const hovered: Array<Rect | null> = Array(count).fill(null);

  const gapFactor = 1 - size[0];
  const maxWidth = (size[1] ?? Infinity) * pxRatio;

  function putBox(
    ctx: CanvasRenderingContext2D,
    rect: uPlot.RectH,
    xOff: number,
    yOff: number,
    lft: number,
    top: number,
    wid: number,
    hgt: number,
    strokeWidth: number,
    iy: number,
    ix: number,
    value: any
  ) {
    if (strokeWidth) {
      ctx.beginPath();
      rect(ctx, lft + strokeWidth / 2, top + strokeWidth / 2, wid - strokeWidth, hgt - strokeWidth);
      ctx.strokeStyle = stroke(iy, ix, value);
      ctx.stroke();
    }

    ctx.beginPath();
    rect(ctx, lft, top, wid, hgt);
    ctx.fillStyle = fill(iy, ix, value);
    ctx.fill();

    qt.add({
      x: round(lft - xOff),
      y: round(top - yOff),
      w: wid,
      h: hgt,
      sidx: iy + 1,
      didx: ix,
    });
  }

  const drawPaths: Series.PathBuilder = (u, sidx, idx0, idx1) => {
    uPlot.orient(
      u,
      sidx,
      (series, dataX, dataY, scaleX, scaleY, valToPosX, valToPosY, xOff, yOff, xDim, yDim, moveTo, lineTo, rect) => {
        let strokeWidth = round((series.width || 0) * pxRatio);

        u.ctx.save();
        rect(u.ctx, u.bbox.left, u.bbox.top, u.bbox.width, u.bbox.height);
        u.ctx.clip();

        walk(laneWidth, sidx - 1, count, yDim, (iy, y0, hgt) => {
          if (mode === TimelineMode.Spans) {
            for (let ix = 0; ix < dataY.length; ix++) {
              if (dataY[ix] != null) {
                let lft = Math.round(valToPosX(dataX[ix], scaleX, xDim, xOff));

                let nextIx = ix;
                while (dataY[++nextIx] === undefined && nextIx < dataY.length) {}

                // to now (not to end of chart)
                let rgt =
                  nextIx === dataY.length
                    ? xOff + xDim + strokeWidth
                    : Math.round(valToPosX(dataX[nextIx], scaleX, xDim, xOff));

                putBox(
                  u.ctx,
                  rect,
                  xOff,
                  yOff,
                  lft,
                  round(yOff + y0),
                  rgt - lft,
                  round(hgt),
                  strokeWidth,
                  iy,
                  ix,
                  dataY[ix]
                );

                ix = nextIx - 1;
              }
            }
          } else if (mode === TimelineMode.Grid) {
            let colWid = valToPosX(dataX[1], scaleX, xDim, xOff) - valToPosX(dataX[0], scaleX, xDim, xOff);
            let gapWid = colWid * gapFactor;
            let barWid = round(min(maxWidth, colWid - gapWid) - strokeWidth);
            let xShift = align === 1 ? 0 : align === -1 ? barWid : barWid / 2;

            for (let ix = idx0; ix <= idx1; ix++) {
              if (dataY[ix] != null) {
                // TODO: all xPos can be pre-computed once for all series in aligned set
                let lft = valToPosX(dataX[ix], scaleX, xDim, xOff);

                putBox(
                  u.ctx,
                  rect,
                  xOff,
                  yOff,
                  round(lft - xShift),
                  round(yOff + y0),
                  barWid,
                  round(hgt),
                  strokeWidth,
                  iy,
                  ix,
                  dataY[ix]
                );
              }
            }
          }
        });

        u.ctx.restore();
      }
    );

    return null;
  };

  const drawPoints: Series.Points.Show =
    formatValue == null
      ? false
      : (u, sidx, i0, i1) => {
          u.ctx.save();
          u.ctx.rect(u.bbox.left, u.bbox.top, u.bbox.width, u.bbox.height);
          u.ctx.clip();

          u.ctx.font = font;
          u.ctx.fillStyle = 'black';
          u.ctx.textAlign = mode === TimelineMode.Spans ? 'left' : 'center';
          u.ctx.textBaseline = 'middle';

          uPlot.orient(
            u,
            sidx,
            (
              series,
              dataX,
              dataY,
              scaleX,
              scaleY,
              valToPosX,
              valToPosY,
              xOff,
              yOff,
              xDim,
              yDim,
              moveTo,
              lineTo,
              rect
            ) => {
              let y = round(yOff + yMids[sidx - 1]);

              for (let ix = 0; ix < dataY.length; ix++) {
                if (dataY[ix] != null) {
                  let x = valToPosX(dataX[ix], scaleX, xDim, xOff);
                  u.ctx.fillText(formatValue(sidx, dataY[ix]), x, y);
                }
              }
            }
          );

          u.ctx.restore();

          return false;
        };

  const init = (u: uPlot) => {
    let over = u.root.querySelector('.u-over')! as HTMLElement;
    over.style.overflow = 'hidden';
    hoverMarks.forEach((m) => {
      over.appendChild(m);
    });
  };

  const drawClear = (u: uPlot) => {
    qt = qt || new Quadtree(0, 0, u.bbox.width, u.bbox.height);

    qt.clear();

    // force-clear the path cache to cause drawBars() to rebuild new quadtree
    u.series.forEach((s) => {
      // @ts-ignore
      s._paths = null;
    });
  };

  const setCursor = (u: uPlot) => {
    let cx = round(u.cursor!.left! * pxRatio);

    for (let i = 0; i < count; i++) {
      let found: Rect | null = null;

      if (cx >= 0) {
        let cy = yMids[i];

        qt.get(cx, cy, 1, 1, (o) => {
          if (pointWithin(cx, cy, o.x, o.y, o.x + o.w, o.y + o.h)) {
            found = o;
          }
        });
      }

      let h = hoverMarks[i];

      if (found) {
        if (found !== hovered[i]) {
          hovered[i] = found;

          h.style.display = '';
          h.style.left = round(found!.x / pxRatio) + 'px';
          h.style.top = round(found!.y / pxRatio) + 'px';
          h.style.width = round(found!.w / pxRatio) + 'px';
          h.style.height = round(found!.h / pxRatio) + 'px';
        }
      } else if (hovered[i] != null) {
        h.style.display = 'none';
        hovered[i] = null;
      }
    }
  };

  // hide y crosshair & hover points
  const cursor: Partial<Cursor> = {
    y: false,
    points: { show: false },
  };

  const yMids: number[] = Array(count).fill(0);
  const ySplits: number[] = Array(count).fill(0);

  return {
    cursor,

    ySplits: (u: uPlot) => {
      walk(laneWidth, null, count, u.bbox.height, (iy, y0, hgt) => {
        // vertical midpoints of each series' timeline (stored relative to .u-over)
        yMids[iy] = round(y0 + hgt / 2);
        ySplits[iy] = u.posToVal(yMids[iy] / pxRatio, FIXED_UNIT);
      });

      return ySplits;
    },
    yValues: (u: uPlot, splits: number[]) => splits.map((v, i) => label(i + 1)),

    yRange: [0, 1] as uPlot.Range.MinMax,

    // pathbuilders
    drawPaths,
    drawPoints,

    // hooks
    init,
    drawClear,
    setCursor,
  };
}
