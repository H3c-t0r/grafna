/**
 * This file holds generic kubernetes compatible types.
 *
 * This is very much a work in progress aiming to simplify common access patterns for k8s resource
 * Please update and improve types/utilities while we find a good pattern here!
 *
 * Once this is more stable and represents a more general pattern, it should be moved to @grafana/data
 *
 */

export interface TypeMeta {
  apiVersion: string;
  kind: string;
}

export interface ObjectMeta {
  // Name is the unique identifier in k8s -- it maps to the "uid" value in most existing grafana objects
  name: string;
  // Namespace maps the owner group -- it is typically the org or stackId for most grafana resources
  namespace?: string;
  // Resource version will increase (not sequentially!) with any change to the saved value
  resourceVersion: string;
  // The first time this was saved
  creationTimestamp: string;
  // General resource annotations -- including the common grafana.app values
  annotations?: GrafanaAnnotations;
  // General application level key+value pairs
  labels?: Record<string, string>;
}

export const AnnoKeyCreatedBy = 'grafana.app/createdBy';
export const AnnoKeyUpdatedTimestamp = 'grafana.app/updatedTimestamp';
export const AnnoKeyUpdatedBy = 'grafana.app/updatedBy';
export const AnnoKeyFolder = 'grafana.app/folder';
export const AnnoKeySlug = 'grafana.app/slug';

// Identify where values came from
const AnnoKeyOriginName = 'grafana.app/originName';
const AnnoKeyOriginPath = 'grafana.app/originPath';
const AnnoKeyOriginKey = 'grafana.app/originKey';
const AnnoKeyOriginTimestamp = 'grafana.app/originTimestamp';

type GrafanaAnnotations = {
  [AnnoKeyCreatedBy]?: string;
  [AnnoKeyUpdatedTimestamp]?: string;
  [AnnoKeyUpdatedBy]?: string;
  [AnnoKeyFolder]?: string;
  [AnnoKeySlug]?: string;

  [AnnoKeyOriginName]?: string;
  [AnnoKeyOriginPath]?: string;
  [AnnoKeyOriginKey]?: string;
  [AnnoKeyOriginTimestamp]?: string;

  // Any key value
  [key: string]: string | undefined;
};

export interface Resource<T = object> extends TypeMeta {
  metadata: ObjectMeta;
  spec: T;
}

export interface ResourceWithStatus<T = object, S = object> extends TypeMeta {
  metadata: ObjectMeta;
  spec: T;
  status: S;
}

export interface ResourceForCreate<T = object> extends Partial<TypeMeta> {
  metadata: Partial<ObjectMeta>;
  spec: T;
}

export interface ListMeta {
  resourceVersion: string;
  continue?: string;
  remainingItemCount?: number;
}

export interface ResourceList<T> extends TypeMeta {
  metadata: ListMeta;
  items: Array<Resource<T>>;
}
