/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2026 emschu[aet]mailbox.org
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public
 * License along with this program.
 * If not, see <https://www.gnu.org/licenses/>.
 */
/**
 * this abstract class provides the readmore click-logic used in recommendations and search section.
 * should be used in combination with the {@link ReadMorePipe }
 */
export abstract class AbstractReadMoreComponent {
  /**
   * call this method from the subclasses
   *
   * @param e an Event instance
   */
  onReadMore(e: any): void {
    if (e.target && e.target.parentNode.classList.contains('keyword-search')) {
      const parentNode: Element = e.target.parentNode;
      // inject a "span"
      parentNode?.parentNode?.append('</span>');
      if (e.target.nextSibling) {
        e.target.nextSibling.innerText = e.target.nextSibling.innerText.replace('</span>');
      }
      parentNode.classList.add('read-more-break');
    } else {
      e.target?.nextSibling?.classList.remove('d-hide');
    }
    e.target?.classList.add('d-hide');
    e.target?.nextSibling?.nextSibling?.classList?.remove('d-hide');
    // TODO add/implement less btn?
  }
}
