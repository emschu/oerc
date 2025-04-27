import { Component } from '@angular/core';
import dayjs from "dayjs";
import {environment} from "../../../environments/environment";

@Component({
  selector: 'app-xmltv',
  templateUrl: './xmltv.component.html',
  styleUrl: './xmltv.component.scss'
})
export class XmltvComponent {
  exampleXmlTvApiUrl: string = "";

  constructor() {
    this.changeExampleXmlTvUrl(XmlTvTimeExpression.TODAY_PLUS_14_DAYS)
  }

  openXmlTv(timeExpression: XmlTvTimeExpression = XmlTvTimeExpression.TODAY_PLUS_7_DAYS) {
    let {from, to} = this.convertTimeExpression(timeExpression);

    window.open(`${environment.serverEndpoint}/xmltv?from=${from.toISOString()}&to=${to.toISOString()}`, '_blank')?.focus()
  }

  private convertTimeExpression(timeExpression: XmlTvTimeExpression) {
    let from = dayjs()
    let to = dayjs().add(1, 'day')

    switch (timeExpression) {
      case XmlTvTimeExpression.TODAY:
        break;
      case XmlTvTimeExpression.TOMORROW:
        from = from.add(1, 'day')
        break;
      case XmlTvTimeExpression.YESTERDAY:
        from = from.subtract(1, 'day')
        break;
      case XmlTvTimeExpression.THIS_WEEK:
        from = from.day(0)
        to = from.day(6)
        break;
      case XmlTvTimeExpression.TODAY_PLUS_7_DAYS:
        to = from.add(7, 'day')
        break;
      case XmlTvTimeExpression.TODAY_PLUS_14_DAYS:
        to = from.add(14, 'day')
        break;
    }
    return {from, to};
  }

  changeExampleXmlTvUrl(timeExpression: XmlTvTimeExpression = XmlTvTimeExpression.TODAY_PLUS_7_DAYS) {
    let {from, to} = this.convertTimeExpression(timeExpression);
    this.exampleXmlTvApiUrl = `${window.location.protocol}//${window.location.host}${environment.serverEndpoint}/xmltv?from=${from.toISOString()}&to=${to.toISOString()}`;
  }

  readonly XmlTvTimeExpression = XmlTvTimeExpression;
}


enum XmlTvTimeExpression {
  TODAY = 'today',
  TOMORROW = 'tomorrow',
  YESTERDAY = 'yesterday',
  THIS_WEEK = 'this_week',
  TODAY_PLUS_7_DAYS = 'today_plus_7_days',
  TODAY_PLUS_14_DAYS = 'today_plus_14_days',
}
