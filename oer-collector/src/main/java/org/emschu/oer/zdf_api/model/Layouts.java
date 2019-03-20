
package org.emschu.oer.zdf_api.model;

/*-
 * #%L
 * oer-server
 * %%
 * Copyright (C) 2019 emschu[aet]mailbox.org
 * %%
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * #L%
 */

import com.google.gson.annotations.Expose;
import com.google.gson.annotations.SerializedName;

import java.util.ArrayList;
import java.util.List;

public class Layouts {

    @SerializedName("1900x200")
    @Expose
    private String _1900x200;
    @SerializedName("1900x400")
    @Expose
    private String _1900x400;
    @SerializedName("640x720")
    @Expose
    private String _640x720;
    @SerializedName("768xauto")
    @Expose
    private String _768xauto;
    @SerializedName("1280xauto")
    @Expose
    private String _1280xauto;
    @SerializedName("240x270")
    @Expose
    private String _240x270;
    @SerializedName("384x216")
    @Expose
    private String _384x216;
    @SerializedName("760x340")
    @Expose
    private String _760x340;
    @SerializedName("840x280")
    @Expose
    private String _840x280;
    @SerializedName("2400x1350")
    @Expose
    private String _2400x1350;
    @SerializedName("1500x300")
    @Expose
    private String _1500x300;
    @SerializedName("840x360")
    @Expose
    private String _840x360;
    @SerializedName("840x140")
    @Expose
    private String _840x140;
    @SerializedName("1900x570")
    @Expose
    private String _1900x570;
    @SerializedName("original")
    @Expose
    private String original;
    @SerializedName("384xauto")
    @Expose
    private String _384xauto;
    @SerializedName("1920x1080")
    @Expose
    private String _1920x1080;
    @SerializedName("276x155")
    @Expose
    private String _276x155;
    @SerializedName("1280x720")
    @Expose
    private String _1280x720;
    @SerializedName("1500x600")
    @Expose
    private String _1500x600;
    @SerializedName("1500x800")
    @Expose
    private String _1500x800;
    @SerializedName("380x170")
    @Expose
    private String _380x170;
    @SerializedName("1152x1296")
    @Expose
    private String _1152x1296;
    @SerializedName("314x314")
    @Expose
    private String _314x314;
    @SerializedName("768x432")
    @Expose
    private String _768x432;

    public List<String> getAllImageLinks() {
        List<String> allLinks = new ArrayList<>();
        allLinks.add(get240x270());
        allLinks.add(get276x155());
        allLinks.add(get314x314());
        allLinks.add(get380x170());
        allLinks.add(get384x216());
        allLinks.add(get640x720());
        allLinks.add(get760x340());
        allLinks.add(get768x432());
        allLinks.add(get840x140());
        allLinks.add(get840x280());
        allLinks.add(get840x360());
        allLinks.add(get1152x1296());
        allLinks.add(get1280x720());
        allLinks.add(get1500x300());
        allLinks.add(get1500x600());
        allLinks.add(get1500x800());
        allLinks.add(get1900x200());
        allLinks.add(get1900x400());
        allLinks.add(get1900x570());
        allLinks.add(get1920x1080());
        allLinks.add(get2400x1350());
        allLinks.add(getOriginal());

        // and now filter out null and empty values
        List<String> imageLinks = new ArrayList<>();
        for (String imageLink : allLinks) {
            if (imageLink != null && !imageLink.isEmpty()) {
                imageLinks.add(imageLink);
            }
        }
        return imageLinks;
    }

    public String get1900x200() {
        return _1900x200;
    }

    public void set1900x200(String _1900x200) {
        this._1900x200 = _1900x200;
    }

    public Layouts with1900x200(String _1900x200) {
        this._1900x200 = _1900x200;
        return this;
    }

    public String get1900x400() {
        return _1900x400;
    }

    public void set1900x400(String _1900x400) {
        this._1900x400 = _1900x400;
    }

    public Layouts with1900x400(String _1900x400) {
        this._1900x400 = _1900x400;
        return this;
    }

    public String get640x720() {
        return _640x720;
    }

    public void set640x720(String _640x720) {
        this._640x720 = _640x720;
    }

    public Layouts with640x720(String _640x720) {
        this._640x720 = _640x720;
        return this;
    }

    public String get768xauto() {
        return _768xauto;
    }

    public void set768xauto(String _768xauto) {
        this._768xauto = _768xauto;
    }

    public Layouts with768xauto(String _768xauto) {
        this._768xauto = _768xauto;
        return this;
    }

    public String get1280xauto() {
        return _1280xauto;
    }

    public void set1280xauto(String _1280xauto) {
        this._1280xauto = _1280xauto;
    }

    public Layouts with1280xauto(String _1280xauto) {
        this._1280xauto = _1280xauto;
        return this;
    }

    public String get240x270() {
        return _240x270;
    }

    public void set240x270(String _240x270) {
        this._240x270 = _240x270;
    }

    public Layouts with240x270(String _240x270) {
        this._240x270 = _240x270;
        return this;
    }

    public String get384x216() {
        return _384x216;
    }

    public void set384x216(String _384x216) {
        this._384x216 = _384x216;
    }

    public Layouts with384x216(String _384x216) {
        this._384x216 = _384x216;
        return this;
    }

    public String get760x340() {
        return _760x340;
    }

    public void set760x340(String _760x340) {
        this._760x340 = _760x340;
    }

    public Layouts with760x340(String _760x340) {
        this._760x340 = _760x340;
        return this;
    }

    public String get840x280() {
        return _840x280;
    }

    public void set840x280(String _840x280) {
        this._840x280 = _840x280;
    }

    public Layouts with840x280(String _840x280) {
        this._840x280 = _840x280;
        return this;
    }

    public String get2400x1350() {
        return _2400x1350;
    }

    public void set2400x1350(String _2400x1350) {
        this._2400x1350 = _2400x1350;
    }

    public Layouts with2400x1350(String _2400x1350) {
        this._2400x1350 = _2400x1350;
        return this;
    }

    public String get1500x300() {
        return _1500x300;
    }

    public void set1500x300(String _1500x300) {
        this._1500x300 = _1500x300;
    }

    public Layouts with1500x300(String _1500x300) {
        this._1500x300 = _1500x300;
        return this;
    }

    public String get840x360() {
        return _840x360;
    }

    public void set840x360(String _840x360) {
        this._840x360 = _840x360;
    }

    public Layouts with840x360(String _840x360) {
        this._840x360 = _840x360;
        return this;
    }

    public String get840x140() {
        return _840x140;
    }

    public void set840x140(String _840x140) {
        this._840x140 = _840x140;
    }

    public Layouts with840x140(String _840x140) {
        this._840x140 = _840x140;
        return this;
    }

    public String get1900x570() {
        return _1900x570;
    }

    public void set1900x570(String _1900x570) {
        this._1900x570 = _1900x570;
    }

    public Layouts with1900x570(String _1900x570) {
        this._1900x570 = _1900x570;
        return this;
    }

    public String getOriginal() {
        return original;
    }

    public void setOriginal(String original) {
        this.original = original;
    }

    public Layouts withOriginal(String original) {
        this.original = original;
        return this;
    }

    public String get384xauto() {
        return _384xauto;
    }

    public void set384xauto(String _384xauto) {
        this._384xauto = _384xauto;
    }

    public Layouts with384xauto(String _384xauto) {
        this._384xauto = _384xauto;
        return this;
    }

    public String get1920x1080() {
        return _1920x1080;
    }

    public void set1920x1080(String _1920x1080) {
        this._1920x1080 = _1920x1080;
    }

    public Layouts with1920x1080(String _1920x1080) {
        this._1920x1080 = _1920x1080;
        return this;
    }

    public String get276x155() {
        return _276x155;
    }

    public void set276x155(String _276x155) {
        this._276x155 = _276x155;
    }

    public Layouts with276x155(String _276x155) {
        this._276x155 = _276x155;
        return this;
    }

    public String get1280x720() {
        return _1280x720;
    }

    public void set1280x720(String _1280x720) {
        this._1280x720 = _1280x720;
    }

    public Layouts with1280x720(String _1280x720) {
        this._1280x720 = _1280x720;
        return this;
    }

    public String get1500x600() {
        return _1500x600;
    }

    public void set1500x600(String _1500x600) {
        this._1500x600 = _1500x600;
    }

    public Layouts with1500x600(String _1500x600) {
        this._1500x600 = _1500x600;
        return this;
    }

    public String get1500x800() {
        return _1500x800;
    }

    public void set1500x800(String _1500x800) {
        this._1500x800 = _1500x800;
    }

    public Layouts with1500x800(String _1500x800) {
        this._1500x800 = _1500x800;
        return this;
    }

    public String get380x170() {
        return _380x170;
    }

    public void set380x170(String _380x170) {
        this._380x170 = _380x170;
    }

    public Layouts with380x170(String _380x170) {
        this._380x170 = _380x170;
        return this;
    }

    public String get1152x1296() {
        return _1152x1296;
    }

    public void set1152x1296(String _1152x1296) {
        this._1152x1296 = _1152x1296;
    }

    public Layouts with1152x1296(String _1152x1296) {
        this._1152x1296 = _1152x1296;
        return this;
    }

    public String get314x314() {
        return _314x314;
    }

    public void set314x314(String _314x314) {
        this._314x314 = _314x314;
    }

    public Layouts with314x314(String _314x314) {
        this._314x314 = _314x314;
        return this;
    }

    public String get768x432() {
        return _768x432;
    }

    public void set768x432(String _768x432) {
        this._768x432 = _768x432;
    }

    public Layouts with768x432(String _768x432) {
        this._768x432 = _768x432;
        return this;
    }

}
