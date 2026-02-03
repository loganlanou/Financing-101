-- +goose Up

-- Learning modules (course categories)
CREATE TABLE IF NOT EXISTS learning_modules (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('basics', 'intermediate', 'advanced')),
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Individual lessons within modules
CREATE TABLE IF NOT EXISTS lessons (
    id TEXT PRIMARY KEY,
    module_id TEXT NOT NULL REFERENCES learning_modules(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    summary TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Glossary of financial terms
CREATE TABLE IF NOT EXISTS glossary_terms (
    id TEXT PRIMARY KEY,
    term TEXT NOT NULL UNIQUE,
    definition TEXT NOT NULL,
    category TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Daily learning tips
CREATE TABLE IF NOT EXISTS learning_tips (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT NOT NULL,
    learn_url TEXT,
    active_date DATE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_lessons_module_id ON lessons(module_id);
CREATE INDEX idx_learning_modules_category ON learning_modules(category);
CREATE INDEX idx_glossary_terms_category ON glossary_terms(category);
CREATE INDEX idx_learning_tips_active_date ON learning_tips(active_date);

-- Seed data: Learning Modules (Basics)
INSERT INTO learning_modules (id, title, description, category, sort_order) VALUES
    ('mod-basics-1', 'What is a Stock?', 'Learn the fundamentals of stock ownership and what it means to own a piece of a company.', 'basics', 1),
    ('mod-basics-2', 'Understanding P/E Ratio', 'Discover how the price-to-earnings ratio helps evaluate stock valuations.', 'basics', 2),
    ('mod-basics-3', 'Market Cap Explained', 'Understand how market capitalization categorizes companies by size.', 'basics', 3),
    ('mod-basics-4', 'What is Volatility?', 'Learn about price fluctuations and why they matter for your investments.', 'basics', 4),
    ('mod-basics-5', 'Risk vs Reward Basics', 'Understand the fundamental relationship between potential gains and potential losses.', 'basics', 5);

-- Seed data: Learning Modules (Intermediate)
INSERT INTO learning_modules (id, title, description, category, sort_order) VALUES
    ('mod-inter-1', 'Reading Financial Statements', 'Learn to interpret balance sheets, income statements, and cash flow statements.', 'intermediate', 1),
    ('mod-inter-2', 'Understanding Market Cycles', 'Explore how markets move through expansion, peak, contraction, and trough phases.', 'intermediate', 2),
    ('mod-inter-3', 'Sector Analysis', 'Learn how to evaluate different market sectors and their characteristics.', 'intermediate', 3),
    ('mod-inter-4', 'What Confidence Levels Mean', 'Understand how AI confidence scores are calculated and their limitations.', 'intermediate', 4),
    ('mod-inter-5', 'Diversification Strategies', 'Learn how spreading investments can help manage portfolio risk.', 'intermediate', 5);

-- Seed data: Learning Modules (Advanced)
INSERT INTO learning_modules (id, title, description, category, sort_order) VALUES
    ('mod-adv-1', 'Technical Analysis Fundamentals', 'Introduction to chart patterns and technical indicators.', 'advanced', 1),
    ('mod-adv-2', 'Options Basics', 'Understanding calls, puts, and basic options strategies.', 'advanced', 2),
    ('mod-adv-3', 'Valuation Methods', 'Compare DCF, comparables, and other valuation approaches.', 'advanced', 3);

-- Seed data: Lessons for "What is a Stock?"
INSERT INTO lessons (id, module_id, title, content, summary, sort_order) VALUES
    ('les-basics-1-1', 'mod-basics-1', 'Stock Ownership Basics',
     'When you buy a stock, you''re purchasing a small piece of ownership in a company. This ownership stake is called equity. As a shareholder, you have a claim on the company''s assets and earnings proportional to the number of shares you own.\n\nFor example, if a company has 1 million shares outstanding and you own 100 shares, you own 0.01% of the company. This ownership gives you certain rights, including voting on important company matters and receiving dividends if the company pays them.',
     'Stocks represent ownership in a company, giving you a proportional claim on assets and earnings.', 1),
    ('les-basics-1-2', 'mod-basics-1', 'Common vs Preferred Stock',
     'There are two main types of stock: common and preferred. Common stock gives you voting rights and the potential for capital appreciation, but dividends are not guaranteed. Preferred stock typically doesn''t come with voting rights, but offers fixed dividends and priority claim on assets if the company is liquidated.\n\nMost individual investors own common stock, which is what people typically mean when they talk about "buying stocks."',
     'Common stock offers voting rights and growth potential; preferred stock provides fixed dividends and priority claims.', 2),
    ('les-basics-1-3', 'mod-basics-1', 'Why Stock Prices Change',
     'Stock prices fluctuate based on supply and demand. When more people want to buy a stock than sell it, the price goes up. When more people want to sell than buy, the price goes down.\n\nMany factors influence this demand: company earnings, economic conditions, interest rates, investor sentiment, and news events. Remember: short-term price movements can be unpredictable, but over the long term, stock prices tend to reflect the underlying business performance.',
     'Stock prices are driven by supply and demand, influenced by company performance, economics, and investor sentiment.', 3);

-- Seed data: Glossary Terms
INSERT INTO glossary_terms (id, term, definition, category) VALUES
    ('glos-1', 'P/E Ratio', 'Price-to-Earnings ratio. A valuation metric calculated by dividing a stock''s price by its earnings per share. A higher P/E may indicate expectations of future growth or overvaluation.', 'Valuation'),
    ('glos-2', 'Market Cap', 'Market Capitalization. The total market value of a company''s outstanding shares, calculated by multiplying share price by total shares. Used to categorize companies as small-cap, mid-cap, or large-cap.', 'Fundamentals'),
    ('glos-3', 'Dividend', 'A portion of company profits distributed to shareholders. Not all companies pay dividends; many growth companies reinvest profits instead.', 'Income'),
    ('glos-4', 'Volatility', 'A measure of how much a stock''s price fluctuates over time. Higher volatility means larger price swings, which indicates higher risk.', 'Risk'),
    ('glos-5', 'Bull Market', 'A market condition where prices are rising or expected to rise. Generally defined as a 20% rise from recent lows.', 'Market Conditions'),
    ('glos-6', 'Bear Market', 'A market condition where prices are falling or expected to fall. Generally defined as a 20% decline from recent highs.', 'Market Conditions'),
    ('glos-7', 'Diversification', 'An investment strategy that spreads money across different assets, sectors, or geographies to reduce risk. The idea is that losses in one area may be offset by gains in another.', 'Strategy'),
    ('glos-8', 'Index', 'A benchmark that tracks the performance of a group of stocks. Examples include the S&P 500 (500 large US companies) and the Dow Jones Industrial Average (30 major companies).', 'Market Conditions'),
    ('glos-9', 'EPS', 'Earnings Per Share. A company''s profit divided by its number of outstanding shares. Used to measure profitability on a per-share basis.', 'Fundamentals'),
    ('glos-10', 'Volume', 'The number of shares traded during a specific period. Higher volume can indicate stronger conviction in price movements.', 'Trading');

-- Seed data: Learning Tips
INSERT INTO learning_tips (id, title, content, category, learn_url) VALUES
    ('tip-1', 'Past Performance Isn''t Everything', 'Historical returns don''t guarantee future results. A stock that''s gone up 50% last year could go down 50% next year. Always look at the fundamentals, not just the chart.', 'Risk', '/learn?module=mod-basics-5'),
    ('tip-2', 'Understand What You Own', 'Before investing in any stock, make sure you understand what the company does, how it makes money, and what risks it faces. If you can''t explain it simply, you might not understand it well enough.', 'Fundamentals', '/learn?module=mod-basics-1'),
    ('tip-3', 'Diversification Matters', 'Don''t put all your eggs in one basket. Spreading investments across different sectors and asset types can help reduce the impact of any single investment''s poor performance.', 'Strategy', '/learn?module=mod-inter-5'),
    ('tip-4', 'AI Insights Are Starting Points', 'Our AI identifies patterns worth investigating, but it cannot predict the future. Use these insights as research prompts, not investment instructions.', 'AI Literacy', '/ai'),
    ('tip-5', 'Emotions and Investing', 'Fear and greed are powerful forces. Many investors buy high (when excited) and sell low (when scared). Having a plan and sticking to it can help manage emotional decisions.', 'Psychology', '/learn?module=mod-basics-5'),
    ('tip-6', 'The Power of Compound Growth', 'Small, consistent returns can grow significantly over time thanks to compounding. A 7% annual return doubles your money roughly every 10 years.', 'Fundamentals', '/learn?module=mod-basics-5'),
    ('tip-7', 'Know Your Risk Tolerance', 'How would you feel if your investment dropped 20%? 50%? Understanding your emotional response to losses helps you choose appropriate investments.', 'Risk', '/learn?module=mod-basics-4');

-- +goose Down
DROP TABLE IF EXISTS learning_tips;
DROP TABLE IF EXISTS glossary_terms;
DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS learning_modules;
