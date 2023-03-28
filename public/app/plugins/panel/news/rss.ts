import { textUtil } from '@grafana/data';

import { getProperty } from './feed';
import { Feed } from './types';

export function parseRSSFeed(txt: string): Feed {
  const domParser = new DOMParser();
  const doc = domParser.parseFromString(textUtil.sanitizeTrustedTypes(txt, 'rss') as unknown as string, 'text/xml');
console.log(doc);
  const feed: Feed = {
    items: Array.from(doc.querySelectorAll('item')).map((node) => ({
      title: getProperty(node, 'title'),
      link: getProperty(node, 'link'),
      content: getProperty(node, 'description'),
      pubDate: getProperty(node, 'pubDate'),
      ogImage: node.querySelector("meta[property='og:image']")?.getAttribute('content'),
    })),
  };

  return feed;
}
