import {Component, OnInit} from '@angular/core';
import {DataSet, Timeline, TimelineOptions, TimelineWindow} from "vis";
import {ApiService} from "../api.service";

@Component({
  selector: 'oer-timeline',
  templateUrl: './timeline.component.html',
  styleUrls: ['./timeline.component.css']
})
export class TimelineComponent implements OnInit {
  private fetchedDays: Date[] = [];
  private timeLine: Timeline;
  private items: DataSet<any> = new DataSet();

  constructor(private apiService: ApiService) {
  }

  ngOnInit() {
    this.initTimeLine();
  }

  initTimeLine() {
    // DOM element where the Timeline will be attached
    const container = document.getElementById('timeline');

    const today: Date = new Date();
    this.fetchDay(today);

    // create groups
    let groups = new DataSet();
    this.apiService.channels().subscribe(value => {
      value.forEach(
        singleChannel => {
          groups.add({id: singleChannel.id, content: singleChannel.name})
        })
    });

    // Configuration for the Timeline
    const options = {
      stack: false,
      start: new Date(),
      end: new Date(1000 * 60 * 60 * 6 + (new Date()).valueOf()),
      editable: false,
      orientation: 'top',
      showCurrentTime: true,
      clickToUse: false,
      horizontalScroll: false,
      verticalScroll: true,
      zoomKey: 'ctrlKey',
      maxHeight: 500,
      moveable: false,
      multiselect: false,
      multiselectPerGroup: false,
      rtl: false,
      selectable: true,
      margin: {
        item: 7,
        axis: 3
      },
      tooltip: {
        followMouse: true,
        overflowMethod: `cap`
      },
    };

    // Create a Timeline
    this.timeLine = new Timeline(container, null, options as TimelineOptions);
    this.timeLine.setGroups(groups);
    this.timeLine.setItems(this.items);

    this.timeLine.on("rangechanged", this.rangeChange.bind(this));
    this.timeLine.on("click", this.itemClicked.bind(this));

    document.getElementById('zoomIn').onclick = this.zoomIn.bind(this);
    document.getElementById('zoomOut').onclick = this.zoomOut.bind(this);
    document.getElementById('moveLeft').onclick = this.moveLeft.bind(this);
    document.getElementById('moveRight').onclick = this.moveRight.bind(this);
  }

  zoomIn(e): void {
    this.timeLine.zoomIn(0.2);
  }
  zoomOut(e): void {
    this.timeLine.zoomOut(0.2);
  }
  moveLeft(e): void {
    this.move(0.2);
  }
  moveRight(e): void {
    this.move(-0.2);
  }

  checkIfDayIsFetched(dayToCheck: Date): boolean {
    let isFetched = false;
    this.fetchedDays.forEach(each => {
      if (each.getDate() == dayToCheck.getDate() && each.getMonth() == dayToCheck.getMonth()
        && each.getFullYear() == dayToCheck.getFullYear()) {
        isFetched = true;
      }
    });
    return isFetched;
  }

  move(percentage: number): void {
    const range: TimelineWindow = this.timeLine.getWindow();
    const interval: number = range.end.valueOf() - range.start.valueOf();

    this.timeLine.setWindow(
      range.start.valueOf() - interval * percentage,
      range.end.valueOf() - interval * percentage
    );
  }

  private async fetchDay(dateToFetch: Date) {
    console.log('fetching day ' + dateToFetch);
    this.fetchedDays.push(new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate()));

    const fromDate: Date = new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate() - 1, 22,0,0,0);
    const toDate: Date = new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate() + 1, 4,0,0,0);

    this.apiService.program(fromDate, toDate).subscribe(value => {
      value.forEach(singleProgramEntry => {
        this.items.update({
          id: singleProgramEntry.id,
          group: singleProgramEntry.channel.id,
          start: singleProgramEntry.startDateTime,
          end: singleProgramEntry.endDateTime,
          content: singleProgramEntry.title,
          title: singleProgramEntry.title,
          type: 'range'
        });
      });
    });
  }

  private rangeChange(e): void {
    if (this.timeLine === undefined) {
      return;
    }
    let rangeStart = this.timeLine.getWindow().start;
    let rangeEnd = this.timeLine.getWindow().end;
    if (!this.checkIfDayIsFetched(rangeStart)) {
      this.fetchDay(new Date(rangeStart.getFullYear(), rangeStart.getMonth(), rangeStart.getDate()));
    } else if (!this.checkIfDayIsFetched(rangeEnd)) {
      this.fetchDay(new Date(rangeEnd.getFullYear(), rangeEnd.getMonth(), rangeEnd.getDate()));
    }
  }

  itemClicked(e): void {
    if (e.item) {
      const clickedEntryId:number = e.item;
      console.log("clicked id:" + clickedEntryId);
      this.apiService.entry(clickedEntryId).subscribe(value => {
        console.log('clicked entry');
        console.log(value);
      });
    }
  }
}
