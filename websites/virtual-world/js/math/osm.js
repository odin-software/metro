const Osm = {
  parseRoads: (data) => {
    const nodes = data.elements.filter((element) => element.type == "node");
    const ways = data.elements.filter((element) => element.type == "way");

    const lats = nodes.map((node) => node.lat);
    const lons = nodes.map((node) => node.lon);

    const minLat = Math.min(...lats);
    const maxLat = Math.max(...lats);
    const minLon = Math.min(...lons);
    const maxLon = Math.max(...lons);

    const deltaLat = maxLat - minLat;
    const deltaLon = maxLon - minLon;
    const aspect = deltaLon / deltaLat;
    const height = deltaLat * 111000 * 10;
    const width = height * aspect * Math.cos(degToRad(maxLat));

    const points = [];
    const segments = [];
    for (const node of nodes) {
      const y = invLerp(maxLat, minLat, node.lat) * height;
      const x = invLerp(minLon, maxLon, node.lon) * width;
      const point = new Point(x, y);
      point.id = node.id;
      points.push(point);
    }
    for (const way of ways) {
      const ids = way.nodes;
      for (let i = 1; i < ids.length; i++) {
        const prev = points.find((p) => p.id == ids[i - 1]);
        const cur = points.find((p) => p.id == ids[i]);
        const oneWay = way.tags && way.tags.oneway || way.tags && way.tags.lanes == 1;

        const segment = new Segment(prev, cur, oneWay);
        segments.push(segment);
      }
    }

    return { points, segments };
  }
}