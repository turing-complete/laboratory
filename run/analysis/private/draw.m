function draw(observed, predicted)
  [~, ~, delta] = kstest2(observed, predicted);

  subplot(2, 2, 1);
  title(sprintf('CDF (delta %.4e)',delta));
  [y, x] = ksdensity(observed, 'function', 'cdf');
  line(x, y, 'LineStyle', '-');
  [y, x] = ksdensity(predicted, 'function', 'cdf');
  line(x, y, 'LineStyle', '--');
  legend('Observed', 'Predicted');

  subplot(2, 2, 2);
  title('PDF');
  [y, x] = ksdensity(observed, 'function', 'pdf');
  line(x, y, 'LineStyle', '-');
  [y, x] = ksdensity(predicted, 'function', 'pdf');
  line(x, y, 'LineStyle', '--');
  legend('Observed', 'Predicted');

  subplot(2, 2, 3);
  hist(observed, 100);
  title('Observed');
  subplot(2, 2, 4);
  hist(predicted, 100);
  title('Predicted');
end
