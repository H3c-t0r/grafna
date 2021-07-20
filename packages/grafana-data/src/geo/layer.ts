import { RegistryItemWithOptions } from '../utils/Registry';
import BaseLayer from 'ol/layer/Base';
import Map from 'ol/Map';
import { PanelData } from '../types';
import { GrafanaTheme2 } from '../themes';
import { PanelOptionsEditorBuilder } from '../utils';
import { ReactNode } from 'react';

/**
 * @alpha
 */
export enum FrameGeometrySourceMode {
  Auto = 'auto', // Will scan fields and find best match
  Geohash = 'geohash',
  Coords = 'coords', // lon field, lat field
  Lookup = 'lookup', // keys > location
  // H3 = 'h3',
  // WKT = 'wkt,
  // geojson? geometry text
}

export enum LookupGeometrySource {
  Countries = 'countries',
  Countries_3Letter = 'countries_3letter',
  Probes = 'probes',
  States = 'states',
}

/**
 * @alpha
 */
export interface FrameGeometrySource {
  mode: FrameGeometrySourceMode;

  // Field mappings
  geohash?: string;
  latitude?: string;
  longitude?: string;
  h3?: string;
  wkt?: string;
  lookup?: string;

  // Path to a mappings file
  lookupSrc?: string;
}

/**
 * This gets saved in panel json
 *
 * depending on the type, it may have additional config
 *
 * This exists in `grafana/data` so the types are well known and extendable but the
 * layout/frame is control by the map panel
 *
 * @alpha
 */
export interface MapLayerOptions<TConfig = any> {
  type: string;
  name?: string; // configured display name

  // Custom options depending on the type
  config?: TConfig;

  // Common method to define geometry fields
  location?: FrameGeometrySource;

  // Common properties:
  // https://openlayers.org/en/latest/apidoc/module-ol_layer_Base-BaseLayer.html
  // Layer opacity (0-1)
  opacity?: number;
}

/**
 * @alpha
 */
export interface MapLayerHandler {
  init: () => BaseLayer;
  legend?: () => ReactNode;
  update?: (data: PanelData) => void;
}

/**
 * Map layer configuration
 *
 * @alpha
 */
export interface MapLayerRegistryItem<TConfig = MapLayerOptions> extends RegistryItemWithOptions {
  /**
   * This layer can be used as a background
   */
  isBaseMap?: boolean;

  /**
   * Show location controls
   */
  showLocation?: boolean;

  /**
   * Show transparency controls in UI (for non-basemaps)
   */
  showOpacity?: boolean;

  /**
   * Function that configures transformation and returns a transformer
   * @param options
   */
  create: (map: Map, options: MapLayerOptions<TConfig>, theme: GrafanaTheme2) => MapLayerHandler;

  /**
   * Show custom elements in the panel edit UI
   */
  registerOptionsUI?: (builder: PanelOptionsEditorBuilder<TConfig>) => void;
}
