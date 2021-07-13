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
