export interface LocationUpdate {
  path?: string;
  query?: UrlQueryMap;
  routeParams?: UrlQueryMap;
  partial?: boolean;
  /*
   * If true this will replace url state (ie cause no new browser history)
   */
  replace?: boolean;
  state?: {
    playlistRunning?: boolean;
  };
}

export interface LocationState {
  url: string;
  path: string;
  query: UrlQueryMap;
  routeParams: UrlQueryMap;
  replace: boolean;
  lastUpdated: number;
  state?: {
    playlistRunning?: boolean;
  };
}

export type UrlQueryValue = string | number | boolean | string[] | number[] | boolean[];
export type UrlQueryMap = { [s: string]: UrlQueryValue };
