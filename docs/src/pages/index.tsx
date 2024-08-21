import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
	const { siteConfig } = useDocusaurusContext();
	return (
		<header className={clsx('hero hero--primary', styles.heroBanner)} style={{
			height: "100vh",
		}}>
			<div className="container">
				<Heading as="h1" className="hero__title">
					{siteConfig.title}
				</Heading>
				<p className="hero__subtitle">{siteConfig.tagline}</p>
				<div className={styles.buttons}>
					<Link
						className="button button--secondary button--lg"
						to="/docs/introduction">
						Docs
					</Link>
				</div>
			</div>
		</header >
	);
}

export default function Home(): JSX.Element {
	const { siteConfig } = useDocusaurusContext();
	return (
		<Layout
			description={`Anvil: ${siteConfig.tagline}`}>
			<HomepageHeader />
		</Layout>
	);
}
