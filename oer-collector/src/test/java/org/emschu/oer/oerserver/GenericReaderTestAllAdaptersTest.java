package org.emschu.oer.oerserver;

/*-
 * #%L
 * oer-collector-project
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

import org.emschu.oer.collector.OerCollector;
import org.emschu.oer.collector.reader.AbstractReader;
import org.emschu.oer.collector.reader.ParserException;
import org.emschu.oer.collector.reader.ProgramEntryPostProcessThread;
import org.emschu.oer.collector.reader.parser.ProgramEntryParserInterface;
import org.emschu.oer.collector.reader.parser.TvShowParserException;
import org.emschu.oer.collector.reader.parser.TvShowParserInterface;
import org.emschu.oer.collector.reader.parser.ard.ARDReader;
import org.emschu.oer.collector.reader.parser.ard.ProgramEntryParser;
import org.emschu.oer.collector.reader.parser.ard.TvShowParser;
import org.emschu.oer.collector.reader.parser.orf.ORFReader;
import org.emschu.oer.collector.reader.parser.srf.SRFReader;
import org.emschu.oer.collector.reader.parser.zdf.ZDFReader;
import org.emschu.oer.collector.service.ChannelService;
import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.ProgramEntry;
import org.emschu.oer.core.model.TvShow;
import org.emschu.oer.core.model.repository.ProgramEntryRepository;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.junit4.SpringRunner;

import java.time.LocalDate;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.stream.StreamSupport;

// TODO add srf reader
@SpringBootTest(classes = OerCollector.class)
@RunWith(SpringRunner.class)
public class GenericReaderTestAllAdaptersTest {

    // adapter reader instances
    @Autowired
    private ARDReader ardReader;

    @Autowired
    private ZDFReader zdfReader;

    @Autowired
    private ORFReader orfReader;

    @Autowired
    private SRFReader srfReader;

    @Autowired
    private ChannelService channelService;

    @Autowired
    private ProgramEntryRepository programEntryRepository;

    @Before
    public void setUp() {
        Assert.assertNotNull(channelService);
        channelService.init();
    }

    @Test
    public void testArdReaderSetup() {
        Assert.assertNotNull(ardReader);
        Assert.assertEquals(Channel.AdapterFamily.ARD, ardReader.getAdapterFamily());
        Assert.assertEquals(ProgramEntryParser.class, ardReader.getProgramEntryParser().getClass());
        Assert.assertEquals(TvShowParser.class, ardReader.getTvShowParser().getClass());

        // check init was called and parsers are registered
        Assert.assertNotNull(ardReader.getCustomParsers());
        Assert.assertEquals(2, ardReader.getCustomParsers().size());
        ardReader.checkConfiguration();
    }

    @Test
    public void testZdfReaderSetup() {
        Assert.assertNotNull(zdfReader);
        Assert.assertEquals(Channel.AdapterFamily.ZDF, zdfReader.getAdapterFamily());
        Assert.assertEquals(org.emschu.oer.collector.reader.parser.zdf.ProgramEntryParser.class, zdfReader.getProgramEntryParser().getClass());
        Assert.assertEquals(org.emschu.oer.collector.reader.parser.zdf.TvShowParser.class, zdfReader.getTvShowParser().getClass());

        // check init was called and parsers are registered
        Assert.assertNotNull(zdfReader.getCustomParsers());
        Assert.assertEquals(0, zdfReader.getCustomParsers().size());
        zdfReader.checkConfiguration();
    }

    @Test
    public void testOrfReaderSetup() {
        Assert.assertNotNull(orfReader);
        Assert.assertEquals(Channel.AdapterFamily.ORF, orfReader.getAdapterFamily());
        Assert.assertEquals(org.emschu.oer.collector.reader.parser.orf.ProgramEntryParser.class, orfReader.getProgramEntryParser().getClass());
        Assert.assertEquals(org.emschu.oer.collector.reader.parser.orf.TvShowParser.class, orfReader.getTvShowParser().getClass());

        // check init was called and parsers are registered
        Assert.assertNotNull(orfReader.getCustomParsers());
        Assert.assertEquals(0, orfReader.getCustomParsers().size());
        orfReader.checkConfiguration();
    }

    @Test
    public void testSrfReaderSetup() {
        Assert.assertNotNull(srfReader);
        Assert.assertEquals(Channel.AdapterFamily.SRF, srfReader.getAdapterFamily());
        Assert.assertEquals(org.emschu.oer.collector.reader.parser.srf.ProgramEntryParser.class, srfReader.getProgramEntryParser().getClass());
        Assert.assertEquals(org.emschu.oer.collector.reader.parser.srf.TvShowParser.class, srfReader.getTvShowParser().getClass());

        // check init was called and parsers are registered
        Assert.assertNotNull(srfReader.getCustomParsers());
        Assert.assertEquals(0, srfReader.getCustomParsers().size());
        srfReader.checkConfiguration();
    }

    @Test
    public void testProgramEntriesAllAdaptersToday() throws ParserException {
        AbstractReader[] abstractReaders = new AbstractReader[]{
                ardReader,
                zdfReader,
                orfReader,
                srfReader
        };
        final LocalDate testDay = LocalDate.now();

        for (AbstractReader abstractReader : abstractReaders) {
            final Channel.AdapterFamily adapterFamily = abstractReader.getAdapterFamily();
            final Channel channel = randomChannel(adapterFamily);

            final ProgramEntryParserInterface programEntryParser1 = abstractReader.getProgramEntryParser();
            final Iterable<?> elements = programEntryParser1.getElements(channel, testDay);
            Assert.assertNotNull(elements);
            final long count = StreamSupport.stream(elements.spliterator(), false).count();
            Assert.assertNotEquals(0, count);
            List<ProgramEntry> testDayProgramEntryList = new ArrayList<>();
            // preprocess
            for (Iterator<?> it = elements.iterator(); it.hasNext(); ) {
                final ProgramEntry programEntry = programEntryParser1.preProcessItem(it.next(), testDay, channel);
                Assert.assertNotNull(programEntry);
                Assert.assertNotNull(programEntry.getTechnicalId());
                if (adapterFamily != Channel.AdapterFamily.ZDF) {
                    Assert.assertNotNull(programEntry.getStartDateTime());
                }
                testDayProgramEntryList.add(programEntry);
                abstractReader.initNewProgramEntry(programEntry, channel);
            }

            programEntryParser1.preProcessProgramList(testDayProgramEntryList);

            // postprocess
            for (ProgramEntry programEntry : testDayProgramEntryList) {
                // workaround for insufficient data of orf parser
                if (adapterFamily == Channel.AdapterFamily.ORF && programEntry.getEndDateTime() == null) {
                    continue;
                }
                new Thread(new ProgramEntryPostProcessThread(programEntryParser1, programEntryRepository, programEntry, true)).run();
                Assert.assertNotNull(programEntry.getId());
                Assert.assertNotNull(programEntry.getChannel());
                Assert.assertEquals(channel.getId(), programEntry.getChannel().getId());
                Assert.assertNotNull(programEntry.getTechnicalId());
                Assert.assertNotNull(programEntry.getCreatedAt());
                Assert.assertNotNull(programEntry.getTitle());
                Assert.assertNotNull(programEntry.getStartDateTime());
                Assert.assertNotNull(programEntry.getEndDateTime());
                Assert.assertNotEquals(0, programEntry.getDurationInMinutes());
                if (adapterFamily != Channel.AdapterFamily.ORF) {
                    Assert.assertNotNull(programEntry.getUrl());
                }
            }

            abstractReader.clear();
        }
    }

    @Test
    public void testTvShowParsers() throws TvShowParserException {
        AbstractReader[] abstractReaders = new AbstractReader[]{
                ardReader,
                zdfReader,
                orfReader
        };
        for (AbstractReader abstractReader : abstractReaders) {
            final TvShowParserInterface tvShowParser = abstractReader.getTvShowParser();
            final Channel.AdapterFamily adapterFamily = abstractReader.getAdapterFamily();

            Assert.assertNotNull(tvShowParser);
            Assert.assertEquals(adapterFamily, tvShowParser.getAdapterFamily());
            List<TvShow> entries = tvShowParser.getEntries();
            Assert.assertNotNull(entries);
            // only work with first ten
            entries = entries.subList(0, 10);

            for (TvShow tvShow : entries) {
                Assert.assertNotNull(tvShow);
                Assert.assertNotNull(tvShow.getTitle());
                tvShowParser.postProcessEntry(tvShow);
                Assert.assertNotNull(tvShow.getAdapterFamily());
                Assert.assertEquals(adapterFamily, tvShow.getAdapterFamily());
                Assert.assertNotNull(tvShow.getUrl());
            }
            tvShowParser.cleanup();
        }
    }

    private Channel randomChannel(Channel.AdapterFamily adapterFamily) {
        final List<Channel> allChannelsByFamily = channelService.getAllChannelsByFamily(adapterFamily);
        Assert.assertNotNull(allChannelsByFamily);
        final int listSize = allChannelsByFamily.size();
        Assert.assertNotEquals(0, listSize);
        return allChannelsByFamily.get((int) (listSize * Math.random()));
    }
}
