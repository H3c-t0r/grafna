import { __awaiter } from "tslib";
import TileLayer from 'ol/layer/Tile';
import OSM from 'ol/source/OSM';
export const standard = {
    id: 'osm-standard',
    name: 'Open Street Map',
    description: 'Add map from a collaborative free geographic world database',
    isBaseMap: true,
    /**
     * Function that configures transformation and returns a transformer
     * @param options
     */
    create: (map, options, eventBus) => __awaiter(void 0, void 0, void 0, function* () {
        return ({
            init: () => {
                return new TileLayer({
                    source: new OSM(),
                });
            },
        });
    }),
};
export const osmLayers = [standard];
//# sourceMappingURL=osm.js.map