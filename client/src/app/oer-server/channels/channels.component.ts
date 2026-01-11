import {Component, inject, OnInit, signal} from '@angular/core';
import {first} from "rxjs/operators";
import {Channel, ChannelResponse} from "../entities";
import {ApiService} from "../api.service";
import {CdkDrag, CdkDragDrop, CdkDropList, moveItemInArray} from "@angular/cdk/drag-drop";

@Component({
  selector: 'app-channels',
  imports: [CdkDrag, CdkDropList],
  templateUrl: './channels.component.html',
  styleUrl: './channels.component.scss',
})
export class ChannelsComponent implements OnInit {

  channels = signal<Channel[]>([]);

  private apiService = inject(ApiService);

  ngOnInit(): void {
    this.apiService.channels().pipe(first()).subscribe((value: ChannelResponse) => {
      this.channels.set(value.data);
    });
  }

  protected drop($event: CdkDragDrop<Channel[]>) {
    const updatedChannels = [...this.channels()];
    moveItemInArray(updatedChannels, $event.previousIndex, $event.currentIndex);
    this.channels.set(updatedChannels);
    this.apiService.updateChannelsOrder(updatedChannels).pipe(first()).subscribe();
  }
}
