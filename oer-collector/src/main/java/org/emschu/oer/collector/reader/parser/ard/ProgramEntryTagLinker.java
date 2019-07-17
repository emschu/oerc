package org.emschu.oer.collector.reader.parser.ard;

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

import org.jsoup.nodes.Document;
import org.jsoup.select.Elements;
import org.emschu.oer.collector.reader.Fetcher;
import org.emschu.oer.collector.reader.ParserException;
import org.emschu.oer.collector.reader.parser.CustomParser;
import org.emschu.oer.collector.service.TagService;
import org.emschu.oer.collector.util.EidExtractor;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.model.Tag;
import org.emschu.oer.core.model.repository.ProgramEntryRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.time.LocalDate;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.logging.Logger;

/**
 * Note: there are two basic page types of ard program web site:
 * 1: "mainTagPages": we have a daily calendar view - very similar to program preview
 * 2: "subTagPages": only show upcoming program entries on one view and the other view
 * is the archive (~4 weeks back from now)
 */
@Component(value = "ardProgramEntryTagLinker")
public class ProgramEntryTagLinker extends CustomParser {

    private static final Logger LOG = Logger.getLogger(ProgramEntryTagLinker.class.getName());

    private static final String ARD_TAG_PAGE_PREFIX = "https://programm.ard.de/TV/Themenschwerpunkte/";

    private Map<String, String> mainTagPages = new HashMap<>();
    private Map<String, String> subTagPages = new HashMap<>();

    @Autowired
    private ProgramEntryRepository programEntryRepository;

    @Autowired
    private TagService tagService;

    public ProgramEntryTagLinker() {
        // init main tag pages
        // we have to fetch them for each day we need data for
        mainTagPages.put("Film", "Film/Alle-Filme/Alle-Filme");
        mainTagPages.put("Dokumentation", "Dokus--Reportagen/Alle-Dokumentationen/Startseite");
        mainTagPages.put("Kultur", "Musik-und-Kultur/Alle-Kultursendungen/Startseite");
        mainTagPages.put("Ratgeber", "Ratgeber-der-ARD/Alle-Ratgeber/Alle-Ratgeber");
        mainTagPages.put("Magazin", "Ratgeber-der-ARD/Magazine/Startseite");
        mainTagPages.put("Serie", "Serien--Soaps/Serien-von-A-bis-Z/Startseite/Serien-von-A-bis-Z");
        mainTagPages.put("Unterhaltung", "Unterhaltung/Alle-Unterhaltungssendungen/Startseite");
        mainTagPages.put("Show/Quiz", "Unterhaltung/Show--Quiz/Startseite");
        mainTagPages.put("Kabarett/Comedy", "Unterhaltung/Kabarett--Comedy/Startseite");
        mainTagPages.put("Zoogeschichten", "Unterhaltung/Zoogeschichten/Startseite");

        // define sub tag pages
        subTagPages.put("Herzgefühl", "Film/Herzgefuehl/Startseite");
        subTagPages.put("Komödie", "Film/Komoedie/Startseite");
        subTagPages.put("Klassiker", "Film/Klassiker/Startseite");
        subTagPages.put("Heimatfilme", "Film/Heimatfilme/Startseite");
        subTagPages.put("Krimi", "Film/Krimi/Startseite");
        subTagPages.put("Tatort", "Film/Tatort/Startseite");
        subTagPages.put("Polizeiruf 110", "Film/Polizeiruf-110/Startseite");
        subTagPages.put("Drama", "Film/Drama/Startseite");
        // disabled 07/2019
//        subTagPages.put("Action/Abenteuer", "Film/Action-und-Abenteuer/Startseite");
        subTagPages.put("Western", "Film/Western/Startseite");
        subTagPages.put("Kurzfilm", "Film/Kurzfilm/Startseite");

        subTagPages.put("Polit-Talkshow", "Politik/Polit-Talkshows/Startseite");
        subTagPages.put("Nachrichten", "Politik/Nachrichten/Startseite");
        subTagPages.put("Aktuelle-Reportage", "Politik/Aktuelle-Reportagen/Startseite");
        subTagPages.put("Polit-Magazine", "Politik/Politmagazine/Startseite");

        subTagPages.put("Geschichte", "Dokus--Reportagen/Geschichte/Startseite");
        subTagPages.put("Kultur-Reportage", "Dokus--Reportagen/Kultur/Startseite");
        subTagPages.put("Tiere", "Dokus--Reportagen/Tiere/Startseite");
        subTagPages.put("Gesundheit", "Dokus--Reportagen/Gesundheit/Startseite");
        subTagPages.put("Umwelt/Natur", "Dokus--Reportagen/Umwelt-und-Natur/Startseite");
        subTagPages.put("Reise-Dokumentation", "Dokus--Reportagen/Reisen/Startseite");
        subTagPages.put("Eisenbahn", "Dokus--Reportagen/Eisenbahn/Startseite");
        subTagPages.put("Wissenschaft", "Dokus--Reportagen/Wissenschaft/Startseite");
        subTagPages.put("Wissensmagazin", "Dokus--Reportagen/Wissensmagazine/Startseite");

        subTagPages.put("Klassik/Oper/Tanz", "Musik-und-Kultur/Klassik-Oper--Tanz/Startseite");
        subTagPages.put("Popkultur", "Musik-und-Kultur/Popkultur/Startseite");
        subTagPages.put("Jazz", "Musik-und-Kultur/Jazz/Startseite");
        subTagPages.put("Literatur", "Musik-und-Kultur/Literatur/Startseite");
        subTagPages.put("Architektur", "Musik-und-Kultur/Architektur/Startseite");
        subTagPages.put("Kultur-Dokumentation", "Musik-und-Kultur/Kultur-Dokumentationen/Startseite");
        subTagPages.put("Kulturmagazine", "Musik-und-Kultur/Kulturmagazine/Startseite");

        subTagPages.put("Heim-/Gartenratgeber", "Ratgeber-der-ARD/Heim-und-Garten/Startseite");
        subTagPages.put("Reiseratgeber", "Ratgeber-der-ARD/Reisen/Startseite");
        subTagPages.put("Gesundheitsratgeber", "Ratgeber-der-ARD/Gesundheit/Startseite");
        subTagPages.put("Natur-/Umweltratgeber", "Ratgeber-der-ARD/Natur-und-Umwelt/Startseite");
        subTagPages.put("Magazin", "Ratgeber-der-ARD/Magazine/Startseite");

        subTagPages.put("Kochen", "Kochen/Alle-Sendungen/Startseite");

        subTagPages.put("Fußball", "Sport/Fussball-im-TV/Startseite");
        subTagPages.put("Sport", "Sport/Alle-Sportsendungen/Startseite");
        subTagPages.put("Sportmagazin", "Sport/Sportmagazine/Startseite");

        subTagPages.put("Soap/Telenovela", "Serien--Soaps/Soaps-und-Telenovelas/Startseite/Startseite"); // correct!
        subTagPages.put("Dokusoap", "Serien--Soaps/Dokusoaps/Startseite");

        subTagPages.put("Show/Quiz", "Unterhaltung/Show--Quiz/Startseite");
        subTagPages.put("Kabarett/Comedy", "Unterhaltung/Kabarett--Comedy/Startseite");
        subTagPages.put("Schlager/Volksmusik", "Unterhaltung/Schlager--Volksmusik/Startseite");
        subTagPages.put("Talkshow", "Unterhaltung/Talkshows/Startseite");
        // disabled 07/2019
//        subTagPages.put("Mit dem Zug", "Unterhaltung/Mit-dem-Zug/Startseite");
        subTagPages.put("Zoogeschichten", "Themenschwerpunkte/Unterhaltung/Zoogeschichten/Startseite");
    }

    @Override
    public void run() throws ParserException {
        if (!isProgramEntryCollectingEnabled() && !isTvShowCollectingEnabled()) {
            LOG.warning("execution stopped due to collection of tv shows and program entries is disabled");
            return;
        }

        // handle main tag pages
        for (LocalDate dateToFetch : getDateRangeList()) {
            mainTagPages.forEach((tagName, tagPath) -> {
                final String url = ARD_TAG_PAGE_PREFIX + tagPath + "?datum=" +
                        dateToFetch.getDayOfMonth() + "." + dateToFetch.getMonthValue() + "." + dateToFetch.getYear() + "&hour=0&ajaxPageLoad=1";

                processEventListPage(tagName, url);
            });
        }

        // handle sub tag pages
        subTagPages.forEach((tagName, tagPath) -> {
            final String previewUrl = ARD_TAG_PAGE_PREFIX + tagPath + "?ajaxPageLoad=1"; // view 1
            final String archiveUrl = previewUrl + "&archiv=1"; // view 2
            processEventListPage(tagName, previewUrl);

            processEventListPage(tagName, archiveUrl);
        });
    }

    private void processEventListPage(String tagName, String url) {
        List<String> technicalIdList = new ArrayList<>();
        LOG.finest("Fetching: " + url);
        Document doc;
        try {
            doc = Fetcher.fetchDocument(url, ".event-list");
        } catch(IllegalStateException ise) {
            LOG.warning("Something went wrong fetching: " + url);
            return;
        }
        Elements elements = doc.select(".event-list ul li[class^=eid]");
        elements.forEach(
                element -> {
                    String classString = element.attr("class");
                    technicalIdList.add(EidExtractor.extractEid(classString));
                }
        );
        if (technicalIdList.isEmpty()) {
            LOG.warning(String.format("No tv shows found for '%s' on page '%s'", tagName, url));
        }

        linkTagToProgramEntries(tagName, technicalIdList);
    }

    @Override
    public void cleanup() {
        tagService.clear();
    }

    private void linkTagToProgramEntries(String tagName, List<String> technicalIdList) {
        Tag linkableTag = tagService.getOrCreateTag(tagName);

        Iterable<ProgramEntry> programEntries = programEntryRepository.getAllByTechnicalIdIsInAndAdapterFamily(technicalIdList, getAdapterFamily());
        int c = 0;
        for (ProgramEntry singleEntry : programEntries) {
            if (!singleEntry.getTags().contains(linkableTag)) {
                singleEntry.getTags().add(linkableTag);
                programEntryRepository.save(singleEntry);
                c += 1;
            }
        }

        LOG.info(String.format("linked tag '%s' to %d program entries", tagName, c));
    }

    public Map<String, String> getMainTagPages() {
        return mainTagPages;
    }

    public Map<String, String> getSubTagPages() {
        return subTagPages;
    }
}
